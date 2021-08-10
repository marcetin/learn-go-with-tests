# Контекст

**[Сав код за ово поглавље можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/context)**

Софтвер често покреће дуготрајне процесе који захтевају много ресурса (често у гороутинама). Ако се радња која је узроковала ово откаже или из неког разлога не успе, морате да зауставите ове процесе на доследан начин кроз своју апликацију.

Ако не управљате овим, ваша брза Го апликација на коју сте толико поносни могла би имати потешкоћа у отклањању грешака у перформансама.

У овом поглављу користићемо `context` пакета како бисмо лакше управљали дуготрајним процесима.

Почећемо са класичним примером веб сервера који када погоди покрене потенцијално дуготрајан процес како би дохватио неке податке како би се вратио у одговору.

Извешћемо сценарио у којем корисник отказује захтев пре него што се подаци могу преузети и побринућемо се да се од процеса одустане.

Поставио сам неки код на срећан пут да бисмо започели. Ево нашег кода сервера.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, store.Fetch())
	}
}
```

Функција `Server` узима `Store` и враћа нам `http.HandlerFunc`. "Store" је дефинисана као:

```go
type Store interface {
	Fetch() string
}
```

Враћена функција позива `store`-ов `Fetch` метод за добијање података и записује их у одговор.

Имамо одговарајући стуб за `Store` који користимо у тесту.

```go
type StubStore struct {
	response string
}

func (s *StubStore) Fetch() string {
	return s.response
}

func TestServer(t *testing.T) {
	data := "hello, world"
	svr := Server(&StubStore{data})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	svr.ServeHTTP(response, request)

	if response.Body.String() != data {
		t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
	}
}
```

Сада када имамо срећан пут, желимо да направимо реалнији сценарио у коме `Store` не може да заврши `Fetch` пре него што корисник откаже захтев.

## Прво напишите тест

Нашем руковаоцу ће бити потребан начин да каже `Store` да откаже посао, па ажурирајте интерфејс.

```go
type Store interface {
	Fetch() string
	Cancel()
}
```

Мораћемо да прилагодимо нашег шпијуна тако да је потребно неко време да вратимо `data` и начин да знамо да је речено да се поништи. Такође ћемо га преименовати у `SpyStore` јер сада посматрамо начин на који се зове. Мораће да дода `Cancel` као метод за имплементацију интерфејса `Store`.

```go
type SpyStore struct {
	response string
	cancelled bool
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}
```

Додајмо нови тест где отказујемо захтев пре 100 милисекунди и проверавамо "Store" да ли се отказује.

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
      data := "hello, world"
      store := &SpyStore{response: data}
      svr := Server(store)

      request := httptest.NewRequest(http.MethodGet, "/", nil)

      cancellingCtx, cancel := context.WithCancel(request.Context())
      time.AfterFunc(5 * time.Millisecond, cancel)
      request = request.WithContext(cancellingCtx)

      response := httptest.NewRecorder()

      svr.ServeHTTP(response, request)

      if !store.cancelled {
          t.Error("store was not told to cancel")
      }
  })
```

Са [Го Блог:Контекст](https://blog.golang.org/context)

> Пакет контекста пружа функције за извођење нових вредности контекста из постојећих. Ове вредности формирају дрво: када се контекст откаже, сви контексти изведени из њега се такође поништавају.

Важно је да изведете свој контекст тако да се отказивања шире по читавом низу позива за дати захтев.

Оно што ми радимо је да из нашег `захтева` изведемо нови` cancellingCtx` који нам враћа функцију `cancel`. Затим заказујемо да се та функција позове за 5 милисекунди коришћењем `time.AfterFunc`. Коначно, користимо овај нови контекст у свом захтеву позивом `request.WithContext`.

## Покушајте да покренете тест

Тест није успео како смо очекивали.

```go
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.00s)
    	context_test.go:62: store was not told to cancel
```

## Напишите довољно кода да прође

Не заборавите да будете дисциплиновани са ТДД -ом. Напишите _минималну_ количину кода како би наш тест прошао.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.Cancel()
		fmt.Fprint(w, store.Fetch())
	}
}
```

Ово чини овај тест пролазним, али не осећа се добро зар не! Сигурно не бисмо требали отказивати `Store` пре него што дохватимо _сваки захтев_.

То што је дисциплиновано, истакло је недостатак у нашим тестовима, што је добра ствар!

Мораћемо да ажурирамо наш тест срећне путање да бисмо потврдили да се неће отказати.

```go
t.Run("returns data from store", func(t *testing.T) {
    data := "hello, world"
    store := &SpyStore{response: data}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }

    if store.cancelled {
        t.Error("it should not have cancelled the store")
    }
})
```

Покрените оба теста и тест срећне путање (happy path) би сада требао бити неуспешан и сада смо приморани да урадимо разумнију примену.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}
	}
}
```

Шта смо урадили овде?

`context` има методу `Done() `која враћа канал којем се шаље сигнал када је контекст" готов "или" отказан ". Желимо да преслушамо тај сигнал и позовемо `store.Cancel` ако га добијемо, али желимо да га занемаримо ако наша `Store` успе да `Fetch` пре њега.

Да бисмо то управљали, покрећемо `Fetch` у гороутини и она ће записати резултат у нови `data` канала. Затим користимо `select` да бисмо ефикасно прешли на два асинхрона процеса, а затим или пишемо одговор или `Cancel`.

## Рефактор

Можемо мало да преобликујемо наш тестни код тако што ћемо утврдити методе за нашег шпијуна

```go
type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Error("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Error("store was told to cancel")
	}
}
```

Не заборавите да унесете `*testing.T` при креирању шпијуна.

```go
func TestServer(t *testing.T) {
	data := "hello, world"

	t.Run("returns data from store", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}

		store.assertWasNotCancelled()
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		request = request.WithContext(cancellingCtx)

		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		store.assertWasCancelled()
	})
}
```

Овај приступ је у реду, али да ли је идиоматичан?

Има ли смисла да се наш веб сервер брине о ручном отказивању `Store`? Шта ако `Store` такође случајно зависи од других спорих процеса? Мораћемо да се уверимо да `Store.Cancel` исправно преноси отказивање свим својим зависним лицима.

Једна од главних тачака `context` је да је то доследан начин нуђења отказивања.

[Из Го документације](https://golang.org/pkg/context/)

> Долазни захтеви према серверу треба да створе Context, а одлазни позиви серверима треба да прихвате Context. Ланац позива функција између њих мора ширити Context, опционо га замењујући изведеним Context-ом креираним помоћу WithCancel, WithDeadline, WithTimeout или WithValue. Када се Context откаже, сви Context-и изведени из њега се такође поништавају.

Из [Go Blog: Context](https://blog.golang.org/context) поново:

> У "Google"-у захтевамо да Го програмери проследе параметар Context као први аргумент свакој функцији на путањи позива између долазних и одлазних захтева. Ово омогућава Го коду који су развили различити тимови да добро сарађује. Омогућава једноставну контролу над временским ограничењима и отказивањем и осигурава да критичне вредности, попут безбедносних података, правилно пролазе Го програме.

(Застаните на тренутак и размислите о последицама сваке функције коју морате послати у контексту и о ергономији тога.)

Осећате се помало нелагодно? Добро. Покушајмо ипак следити тај приступ и уместо тога проћи кроз `context` у нашу` Store` и пустити га да буде одговоран. На тај начин такође може пренети `context` до својих зависних особа, а и они могу бити одговорни за заустављање.

## Прво напишите тест

Мораћемо да променимо постојеће тестове како се мењају њихове одговорности. Једина ствар за коју је наш руковалац сада одговоран је да пошаље контекст низводној `Store` и да рукује грешком која ће доћи из` Store`-а када се откаже.

Хајде да ажурирамо интерфејс `Store` да бисмо приказали нове одговорности.

```go
type Store interface {
	Fetch(ctx context.Context) (string, error)
}
```

Избришите код унутар нашег руковаоца за сада

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
```

Ажурирајте нашу `SpyStore`

```go
type SpyStore struct {
	response string
	t        *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <-ctx.Done():
				s.t.Log("spy store got cancelled")
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}
}
```

Морамо учинити да се наш шпијун понаша као прави метод који ради са `context`.

Симулирамо спор процес у којем полако градимо резултат додавањем низа, знак по знак у гороутини. Када гороутина заврши свој рад, уписује стринг у `data` канал. Рутина слуша `ctx.Done` и зауставиће рад ако се сигнал пошаље на тај канал.

Коначно, код користи други `select` да сачека да та рутина заврши свој рад или да дође до отказивања.

Слично је нашем приступу од раније, користимо Го-ове истовремене примитиве како бисмо учинили да се два асинхрона процеса међусобно утркују како би одредили шта враћамо.

Сличан приступ ћете заузети приликом писања сопствених функција и метода које прихватају `context`, па се уверите да разумете шта се дешава.

Коначно можемо ажурирати наше тестове. Коментирајте наш тест отказивања како бисмо прво поправили тест сретне стазе.

```go
t.Run("returns data from store", func(t *testing.T) {
    data := "hello, world"
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }
})
```

## Покушајте да покренете тест

```
=== RUN   TestServer/returns_data_from_store
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/returns_data_from_store (0.00s)
    	context_test.go:22: got "", want "hello, world"
```

## Напишите довољно кода да прође

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, _ := store.Fetch(r.Context())
		fmt.Fprint(w, data)
	}
}
```

Наш срећан пут би требао бити ... срећан. Сада можемо поправити други тест.

## Прво напишите тест

Морамо да тестирамо да не напишемо никакав одговор на случај грешке. Нажалост, `httptest.ResponseRecorder` нема начина да то схвати па ћемо морати да откотрљамо сопственог шпијуна да бисмо то тестирали.

```go
type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}
```

Наш `SpyResponseWriter` имплементира `http.ResponseWriter` тако да га можемо користити у тесту.

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)

    cancellingCtx, cancel := context.WithCancel(request.Context())
    time.AfterFunc(5*time.Millisecond, cancel)
    request = request.WithContext(cancellingCtx)

    response := &SpyResponseWriter{}

    svr.ServeHTTP(response, request)

    if response.written {
        t.Error("a response should not have been written")
    }
})
```

## Покушајте да покренете тест

```
=== RUN   TestServer
=== RUN   TestServer/tells_store_to_cancel_work_if_request_is_cancelled
--- FAIL: TestServer (0.01s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.01s)
    	context_test.go:47: a response should not have been written
```

## Напишите довољно кода да прође

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Fetch(r.Context())

		if err != nil {
			return // todo: log error however you like
		}

		fmt.Fprint(w, data)
	}
}
```

Након овога можемо видети да је код сервера поједностављен јер више није експлицитно одговоран за отказивање, већ једноставно пролази кроз `context` и ослања се на функције низводно да поштују сва отказивања која се могу догодити.

## Окончање

463 / 5000
Резултати превода
### Шта смо покрили

- Како тестирати HTTP руковаоца којем је клијент отказао захтев.
- Како користити контекст за управљање отказивањем.
- Како написати функцију која прихвата `context` и користи је за отказивање помоћу гороутина, `select` и канала.
- Следите Гоогле-ове смернице о томе како управљати отказивањем ширењем контекста обухваћеног захтевима кроз ваш низ позива.
- Како пронаћи свог шпијуна за `http.ResponseWriter` ако вам затреба.

### Шта је са context.Value?

[Michal Štrba](https://faiface.github.io/post/context-should-go-away-go2/) и ја имам слично мишљење.

> Ако користите ctx.Value у мојој (непостојећој) компанији, отпуштен си

Неки инжењери су заговарали преношење вредности кроз `context` јер се _осећа згодно_.

Погодност је често узрок лошег кода.

Проблем са `context.Values` је у томе што је то само нетипизирана мапа тако да немате безбедност типова и морате да рукујете њоме која заправо не садржи вашу вредност. Морате створити спрезање кључева карте из једног модула у други и ако неко нешто промени ствари почињу да се кваре.

Укратко, **ако су функцији потребне неке вредности, ставите их као откуцане параметре уместо да их покушавате позвати из `context.Values`**. То га чини статички провереним и документованим да га сви виде.

#### Али...

С друге стране, може бити корисно укључити информације које су ортогоналне према захтеву у контекст, као што је ИД праћења. Потенцијално ове информације не би биле потребне свакој функцији у вашем стогу позива и учиниле би ваше функционалне потписе врло неуредним.
[Jack Lindamood каже **Context.Value треба да обавештава, а не да контролише **](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)

> Садржај context.Value је за одржавање, а не за кориснике. Никада не би требало захтевати улаз за документоване или очекиване резултате.

### Додатни материјал

- Заиста сам уживао читајући [Контекст би требало да нестане у Го 2 аутора Michal Štrba](https://faiface.github.io/post/context-should-go-away-go2/). Његов аргумент је да је свуда преношење `context` мирис, да то указује на недостатак језика у погледу отказивања. Каже да би било боље да се то некако ријеши на нивоу језика, а не на нивоу библиотеке. Док се то не догоди, биће вам потребан `context` ако желите да управљате дуготрајним процесима.
- [Го блог даље описује мотивацију за рад са `context` и има неке примере] (https://blog.golang.org/context)
