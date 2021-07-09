# Врсте грешака

**[Сав код можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/q-and-a/error-types)**

** Стварање сопствених типова за грешке може бити елегантан начин сређивања кода, што олакшава употребу и тестирање кода.**

Пита Педро на `Gopher Slack`-u

> Ако стварам грешку попут `fmt.Errorf("%s must be foo, got %s", bar, baz)`, постоји ли начин да се тестира једнакост без упоређивања вредности низа?

Направимо функцију која ће помоћи у истраживању ове идеје.

```go
// DumbGetter will get the string body of url if it gets a 200
func DumbGetter(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("problem fetching from %s, %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("did not get 200 from %s, got %d", url, res.StatusCode)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body) // ignoring err for brevity

	return string(body), nil
}
```

Неријетко је написати функцију која би могла пропасти из различитих разлога и желимо бити сигурни да правилно поступамо са сваким сценаријем.

Као што Педро каже, `могли` бисмо написати тест за статусну грешку на тај начин.

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	want := fmt.Sprintf("did not get 200 from %s, got %d", svr.URL, http.StatusTeapot)
	got := err.Error()

	if got != want {
		t.Errorf(`got "%v", want "%v"`, got, want)
	}
})
```

Овим тестом се креира сервер који увек враћа `StatusTeapot`, а затим користимо његову УРЛ адресу као аргумент за` DumbGetter` како бисмо могли видети да правилно обрађује одговоре који нису `200`.

## Проблеми са овим начином тестирања

Ова књига покушава да нагласи да _слушате своје тестове_ а овај тест не делује добро:

- Конструишемо исти низ као и производни код да бисмо га тестирали
- Досадно је читати и писати
- Да ли је тачан низ порука о грешци оно што нас заправо занима?

Шта нам ово говори? Ергономија нашег теста одразила би се на још један бит кода који покушава да користи наш код.

Како корисник нашег кода реагује на одређену врсту грешака које враћамо? Најбоље што могу је погледати низ грешака који је изузетно склон грешкама и ужасан за писање.

## Шта треба да радимо

Са ТДД-ом имамо предност уласка у начин размишљања:

> Како бих _ја_ желео да користим овај код?

Оно што бисмо могли учинити за `DumbGetter` је пружање начина да корисници користе систем типа да би разумели каква се грешка догодила.

Шта ако би нам `DumbGetter` могао вратити нешто слично

```go
type BadStatusError struct {
	URL    string
	Status int
}
```

Уместо магичне жице, имамо стварне _податке_ за рад.

Променимо постојећи тест како би одражавао ову потребу

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	got, isStatusErr := err.(BadStatusError)

	if !isStatusErr {
		t.Fatalf("was not a BadStatusError, got %T", err)
	}

	want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
})
```

Мораћемо да натерамо `BadStatusError` да примени интерфејс за грешке.

```go
func (b BadStatusError) Error() string {
	return fmt.Sprintf("did not get 200 from %s, got %d", b.URL, b.Status)
}
```

### Шта ради тест?

Уместо да проверимо тачан низ грешке, ми радимо [тврдњу типа](https://tour.golang.org/methods/15) на грешци како бисмо видели да ли је `BadStatusError`. Ово одражава нашу жељу да _ врста_ грешака буде јаснија. Под претпоставком да тврдња пролази, можемо тада проверити да ли су својства грешке тачна.

Када покренемо тест, говори нам да нисмо вратили праву врсту грешке

```
--- FAIL: TestDumbGetter (0.00s)
    --- FAIL: TestDumbGetter/when_you_dont_get_a_200_you_get_a_status_error (0.00s)
    	error-types_test.go:56: was not a BadStatusError, got *errors.errorString
```

Поправимо `DumbGetter` ажурирањем кода за руковање грешкама како бисмо користили наш тип

```go
if res.StatusCode != http.StatusOK {
	return "", BadStatusError{URL: url, Status: res.StatusCode}
}
```

Ова промена је имала неке _ стварне позитивне ефекте_

- Наша функција `DumbGetter` постала је једноставнија, више се не бави замршеношћу низа грешака, већ само ствара` BadStatusError`.
- Наши тестови сада одражавају (и документују) шта би корисник нашег кода _ могао_ да уради ако одлучи да жели да уради софистицираније поступање са грешкама него само пријављивање. Само урадите тврдњу о типу и тада ћете добити лак приступ својствима грешке.
- То је и даље „само“ `error`, па ако је одлуче могу је проследити низом позива или пријавити као било коју другу `error`.

## Окончање

Ако се нађете на тестирању вишеструких услова грешке, немојте пасти у замку упоређивања порука о грешкама.

То доводи до неуобичајених и тешких тестова за читање / писање и одражава потешкоће које ће имати корисници вашег кода ако и они треба да почну да раде ствари другачије у зависности од врсте грешака које су се догодиле.

Увек се побрините да ваши тестови одражавају како _ бисте_ желели да користите свој код, па у том погледу размислите о стварању типова грешака који ће садржати ваше врсте грешака. Ово олакшава руковање различитим врстама грешака корисницима вашег кода, а такође олакшава и читање писања кода за руковање грешкама.

## Додатак

Од Го 1.13 постоје нови начини за рад са грешкама у стандардној библиотеци која је обрађена на [Го Блог](https://blog.golang.org/go1.13-errors)

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	var got BadStatusError
	isBadStatusError := errors.As(err, &got)
	want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

	if !isBadStatusError {
		t.Fatalf("was not a BadStatusError, got %T", err)
	}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
})
```

У овом случају користимо [`errors.As`](https://golang.org/pkg/errors/#example_As) да бисмо покушали да издвојимо нашу грешку у наш прилагођени тип. Враћа `bool` да означава успех и издваја га у `got` за нас.
