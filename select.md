# Селектовање

**[Сав код за ово поглавље можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/select)**

Од вас је затражено да направите функцију под називом `WebsiteRacer` која узима два УРЛ -а и "трка" их тако што их погађа HTTP GET-ом и враћа URL који се први вратио. Ако се ниједан од њих не врати у року од 10 секунди, требало би да врати `error`.

За ово ћемо користити

- `net/http` за упућивање HTTP позива.
- `net/http/httptest` да нам помогне да их тестирамо.
- Го рутине.
- `select` за синхронизацију процеса.

## Прво напишите тест

Почнимо са нечим простим.

```go
func TestRacer(t *testing.T) {
	slowURL := "http://www.facebook.com"
	fastURL := "http://www.quii.co.uk"

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

Знамо да ово није савршено и да има проблема, али то ће нас покренути. Важно је да се не заморите превише око тога да ствари буду савршене први пут.

## Покушајте да покренете тест

`./racer_test.go:14:9: undefined: Racer`

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

```go
func Racer(a, b string) (winner string) {
	return
}
```

`racer_test.go:25: got '', want 'http://www.quii.co.uk'`

## Напишите довољно кода да прође

```go
func Racer(a, b string) (winner string) {
	startA := time.Now()
	http.Get(a)
	aDuration := time.Since(startA)

	startB := time.Now()
	http.Get(b)
	bDuration := time.Since(startB)

	if aDuration < bDuration {
		return a
	}

	return b
}
```

За сваки URL:

1. Користимо `time.Now()` за снимање непосредно пре него што покушамо да добијемо `URL`.
2. Затим користимо [`http.Get`](https://golang.org/pkg/net/http/#Client.Get) да покушамо да добијемо садржај `URL`. Ова функција враћа [`http.Response`](https://golang.org/pkg/net/http/#Response) као и `error`, али засад нас не занимају ове вредности.
3. `time.Since` узима почетно време и враћа `time.Duration` разлике.

Када ово учинимо, једноставно упоредимо трајање да видимо који је најбржи.

### Проблеми

Ово вам може, али и не мора, проћи тест. Проблем је у томе што се обраћамо правим веб страницама како бисмо тестирали сопствену логику.

Тестирање кода који користи HTTP је толико уобичајено да Го има алате у стандардној библиотеци који ће вам помоћи да га тестирате.

У поглављима исмевања и убризгавања зависности говорили смо о томе како идеално не желимо да се ослањамо на спољне услуге за тестирање нашег кода јер они могу бити

- Спор
- Пахуљице
- Не могу да тестирам ивице

У стандардној библиотеци постоји пакет под називом [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/) где можете лако да направите лажни HTTP сервер.

Хајде да променимо наше тестове тако да користе лажне, тако да имамо поуздане сервере за тестирање против оних које можемо контролисати.

```go
func TestRacer(t *testing.T) {

	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	slowServer.Close()
	fastServer.Close()
}
```

Синтакса може изгледати мало заузето, али само одвојите време.

`httptest.NewServer` узима `http.HandlerFunc` који шаљемо путем _анонимне функције_.

`http.HandlerFunc` је тип који изгледа овако: `type HandlerFunc func(ResponseWriter, *Request)`.

Све што заиста каже је да му је потребна функција која узима `ResponseWriter` и` Request`, што није превише изненађујуће за HTTP сервер.

Испоставило се да овде заиста нема додатне магије, **тако бисте и ви написали _реални_ HTTP сервер у Го **. Једина разлика је у томе што га омотавамо у `httptest.NewServer` што га чини лакшим за тестирање, јер налази отворен порт за слушање, а затим га можете затворити када завршите са тестом.

Унутар наша два сервера, спори ћемо имати кратко `time.Sleep` када добијемо захтев да буде спорији од другог. Оба сервера затим уписују `OK` одговор са `w.WriteHeader(http.StatusOK)` назад позиваоцу.

Ако поново покренете тест, дефинитивно ће проћи сада и требало би да буде бржи. Играјте се са овим "спавањем" да бисте намерно прекинули тест.

## Рефактор

Имамо дуплирање и у нашем производном коду и у тестном коду.

```go
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
```

Ово "СУШЕЊЕ" чини наш `Racer` код много лакшим за читање.

```go
func TestRacer(t *testing.T) {

	slowServer := makeDelayedServer(20 * time.Millisecond)
	fastServer := makeDelayedServer(0 * time.Millisecond)

	defer slowServer.Close()
	defer fastServer.Close()

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
```

Преобликовали смо креирање наших лажних сервера у функцију звану `makeDelayedServer` како бисмо избацили неки незанимљив код из теста и смањили понављање.

### `defer`

Префиксом позива функције са `defer' сада ће позвати ту функцију _на крају функције која садржи_.

Понекад ћете морати да очистите ресурсе, као што је затварање датотеке или у нашем случају затварање сервера тако да не наставља да слуша порт.

Желите да се ово изврши на крају функције, али држите упутство близу места где сте креирали сервер за добробит будућих читалаца кода.

Наше прерађивање је побољшање и разумно је решење с обзиром на Го функције које смо до сада покривали, али решење можемо учинити једноставнијим.

### Синхронизовање процеса

- Зашто тестирамо брзине веб локација једну за другом када је Го одличан у истовремености? Требало би да можемо да проверимо обоје истовремено.
- Не занима нас _тачно време одговора_ на захтеве, само желимо да знамо који ће се први вратити.

Да бисмо то урадили, представићемо нову конструкцију под називом `select` која нам помаже да синхронизујемо процесе заиста лако и јасно.

```go
func Racer(a, b string) (winner string) {
	select {
	case <-ping(a):
		return a
	case <-ping(b):
		return b
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}
```

#### `ping`

Дефинирали смо функцију `ping` која ствара `chan struct{}` и враћа је.

У нашем случају, не _маримо_ који тип се шаље на канал, већ _само желимо да сигнализирамо да смо готови_ и затварање канала ради савршено!

Зашто `struct{}` а не неки други тип попут `bool`? Па, `chan struct{}` је најмањи тип података доступан из меморијске перспективе па смо
не добијају алокацију у односу на `bool`. Пошто затварамо и не шаљемо ништа на каналу, зашто бисмо било шта додељивали?

Унутар исте функције покрећемо рутину која ће слати сигнал у тај канал када завршимо `http.Get(url)`.

##### Увек `make` канале

Обратите пажњу на то како морамо да користимо `make` при креирању канала; уместо да кажете `var ch chan struct{}`. Када користите `var` променљива ће бити иницијализована са" нултом "вредношћу типа. Дакле, за `string` је `"" `, `int` је 0 итд.

За канале нулта вредност је `nil` и ако покушате да јој пошаљете са` <-` блокираће се заувек јер не можете да шаљете на `nil` канале

[Ово можете видети на делу на Го Игралишту](https://play.golang.org/p/IIbeAox5jKA)

#### `select`

Ако се сећате из поглавља о паралелности, можете чекати да се вредности пошаљу на канал са `myVar := <-ch`. Ово је _блокирајући_ позив, јер чекате вредност.

Оно што вам `select` омогућава је да чекате на _мултипле_ канала. Први који пошаље вредност "вин" и код испод `case` се извршава.

Користимо `ping` у свом `select` за постављање два канала за сваки од `URL`-ова. Ко год прво упише на свој канал, његов код ће бити извршен у `select`, што резултира враћањем `URL`-а (и победом).

Након ових промена, намера нашег кода је врло јасна, а имплементација је заправо једноставнија.

### Временска ограничења

Наш последњи захтев био је да вратимо грешку ако `Racer` потраје дуже од 10 секунди.

## Прво напишите тест

```go
t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
	serverA := makeDelayedServer(11 * time.Second)
	serverB := makeDelayedServer(12 * time.Second)

	defer serverA.Close()
	defer serverB.Close()

	_, err := Racer(serverA.URL, serverB.URL)

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
})
```

Учинили смо да нашим тестним серверима треба више од 10 секунди да се врате у овај сценарио и очекујемо да ће `Racer` сада вратити две вредности, победничку УРЛ адресу (коју у овом тесту занемарујемо са `_`) и `error`.

## Покушајте да покренете тест

`./racer_test.go:37:10: assignment mismatch: 2 variables but 1 values`

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	}
}
```

Промените потпис `Racer` да бисте вратили победника и `error`. Вратите `nil` за наше срећне случајеве.

Компајлер ће се жалити на ваш _први тест_ само у потрази за једном вредношћу, па промените ту линију у `got, _ := Racer(slowURL, fastURL)`, знајући да треба да проверимо да _не_ добија грешку у нашем срећном сценарију.

Ако га покренете сада након 11 секунди, неће успети.

```
--- FAIL: TestRacer (12.00s)
    --- FAIL: TestRacer/returns_an_error_if_a_server_doesn't_respond_within_10s (12.00s)
        racer_test.go:40: expected an error but didn't get one
```

## Напишите довољно кода да прође

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

`time.After` је врло згодна функција када се користи `select`. Иако се то у нашем случају није догодило, потенцијално можете написати код који заувијек блокира ако канали на којима слушате никада не врате вриједност. `time.After` враћа `chan` (попут `ping`) и послаће сигнал низ њега након одређеног времена.

За нас је ово савршено; ако `a` или `b` успеју да врате, они побеђују, али ако дођемо до 10 секунди онда ће наше `time.After` послати сигнал, а ми ћемо вратити `error`.

### Спори тестови

Проблем који имамо је што овај тест траје 10 секунди. За тако једноставну логику, ово се не осећа сјајно.

Оно што можемо учинити је да подесимо временско ограничење. Дакле, у нашем тесту можемо имати врло кратко временско ограничење, а онда када се код користи у стварном свету може се поставити на 10 секунди.


```go
func Racer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

Наши тестови се сада неће компајлирати јер не испоручујемо временско ограничење.

Пре него што пожуримо да додамо ову подразумевану вредност у оба теста, _послушајмо их_.

- Да ли нам је стало до истека рока у „срећном“ тесту?
- Захтеви су били експлицитни у вези са временским ограничењем.

С обзиром на ово знање, учинимо мало преобликовање како бисмо били наклоњени и нашим тестовима и корисницима нашег кода.

```go
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

Наши корисници и наш први тест могу да користе `Racer` (који користи `ConfigurableRacer` испод хаубе), а наш тест тужне путање може да користи `ConfigurableRacer`.

```go
func TestRacer(t *testing.T) {

	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within the specified time", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}
```

Додао сам последњу проверу на првом тесту како бих потврдио да не добијамо `error`.

## Окончање

### `select`

- Помаже вам да чекате на више канала.
- Понекад ћете желети да уврстите `time.After` у један од својих `cases` како бисте спречили заувек блокирање система.

### `httptest`

- Погодан начин за креирање тест сервера како бисте имали поуздане и контролисане тестове.
- Коришћење истих интерфејса као и "стварних" `net/http` сервера што је доследно и мање за учење.
