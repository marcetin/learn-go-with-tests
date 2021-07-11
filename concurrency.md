# Конкуренција

**[Сав код за ово поглавље можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/concurrency)**

Ево подешавања: колега је написао функцију `CheckWebsites` проверава статус листе УРЛ адреса.

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = wc(url)
	}

	return results
}
```

Враћа мапу сваког провереног УРЛ-а у логичку вредност - `true` за добро одговор, `false` за лош одговор.

Такође морате проследити `WebsiteChecker` који узима један УРЛ и враћа се боолеан. Ову функцију користи за проверу свих веб локација.

Коришћење [убризгавања пакета од којих зависи апликација][ДИ] им је омогућило да тестирају функцију без упућивање стварних ХТТП позива, чинећи га поузданим и брзим.

Ево теста који су написали:

```go
package concurrency

import (
	"reflect"
	"testing"
)

func mockWebsiteChecker(url string) bool {
	if url == "waat://furhurterwe.geds" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}
```

Функција је у изради и користи се за проверу стотина веб локација. Али ваш колега је почео да добија притужбе да је спор, па су тражили ти да помогнеш да га убрзаш.

## Напишите тест

Користимо референтну вредност да тестирамо брзину `CheckWebsites` како бисмо могли да видимо ефекат наших промена.

```go
package concurrency

import (
	"testing"
	"time"
)

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}
```

Референтни тестови тестирају `CheckWebsites` користећи одсек од стотину УРЛ адреса и употреба нова лажна примена `WebsiteChecker`. `slowStubWebsiteChecker` је намерно споро. Користи `time.Sleep` да сачека тачно двадесет милисекунди и онда се враћа истина.

Када покренемо референтну вредност помоћу `go test -bench=.` (или ако сте у Windows Powershell `go test -bench="."`):

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v0
BenchmarkCheckWebsites-4               1        2249228637 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v0        2.268s
```

`CheckWebsites` је постављен на 2249228637 наносекунди - око две и четврт секунде.

Покушајмо и учинимо ово бржим.

### Напишите довољно кода да прође

Сада коначно можемо разговарати о паралелности која, за потребе следеће значи „имати више ствари у току“. Ово је нешто које природно радимо свакодневно.

На пример, јутрос сам скувао шољу чаја. Ставио сам чајник и онда, док сам чекао да прокључа, извадио сам млеко из фрижидера, добио чај из ормана, пронашао моју омиљену шољу, ставио врећицу чаја у шољу и затим, кад је чајник прокључао, ставио сам воду у шољу.

Оно што нисам _ ставио је чајник и затим стајао ту и буљећи у њега чајник док није прокључао, а затим све остало урадите након што је котлић прокључао.

Ако разумете зашто је брже правити чај на први начин, онда то можете схватите како ћемо убрзати „проверу веб локација“. Уместо да чека веб локација да одговори пре него што пошаље захтев на следећу веб локацију, рећи ћемо наш рачунар да упути следећи захтев док чека.

Обично у Го када функцију зовемо `doSomething()` чекамо да се врати (чак и ако нема вредност за враћање, и даље чекамо да се заврши). Ми то кажемо ова операција *блокира* - чини нас да чекамо да се заврши. Операција који се не блокира у програму Го, извршиће се у одвојеном *процесу* који се назива *гороутина*.
Замислите процес као читање странице Го кода од врха до дна, даље 'унутар' сваке функције када је позову да прочита шта ради. Кад одвојена процес започиње као да други читач почиње читати унутар функције, остављајући оригиналном читачу да настави спуштање странице.

Да бисмо рекли Го-у да започне нови програм, претварамо позив функције у `go` изјава стављањем кључне речи `go` испред ње: `go doSomething() `.

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()
	}

	return results
}
```

Јер једини начин за покретање гороутине је стављање `go` испред функције позива, често користимо *анонимне функције* када желимо да покренемо гороутину. Ан анонимни функцијски литерал изгледа исто као нормална декларација функције, али без имена (што није изненађујуће). Можете да видите један горе у телу петља `for`.

Анонимне функције имају бројне функције које их чине корисним, две од њих које користимо горе. Прво, они могу бити извршени истовремено они су декларисани - то је оно што је `()` на крају анонимне функције радиш. Друго, они задржавају приступ лексичком опсегу у којем су дефинисани - све променљиве које су доступне у тренутку када прогласите анонимним функције су такође доступне у телу функције.

Тело горње анонимне функције је потпуно исто као и тело петље пре него што. Једина разлика је у томе што ће свака итерација петље започети нову гороутина, истовремено са тренутним процесом (функција `WebsiteChecker`) од којих ће сваки додати свој резултат на мапу резултата.

Али када покренемо `go test`:

```sh
--- FAIL: TestCheckWebsites (0.00s)
        CheckWebsites_test.go:31: Wanted map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], got map[]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s

```

### Брза страна у паралелни (изм) универзум ...

Можда нећете добити овај резултат. Можда ћете добити поруку панике мало ћемо разговарати о томе. Не брините ако то имате, само задржите трчање теста док _не_ добијете горњи резултат. Или се претварај да јеси. На вама. Добродошли у паралелност: када се не рукује исправно, тешко је предвидети шта ће се догодити. Не брините - зато пишемо тестове помозите нам да знамо када предвидљиво радимо са паралелношћу.

### ... и вратили смо се.

Ухватили су нас оригинални тестови, `CheckWebsites` сада враћа
празна карта. Шта је пошло наопако?

Ниједна гороутина коју је започела наша петља `for` није имала довољно времена за додавање њихов резултат на мапу `results`; функција `WebsiteChecker` је пребрза за њих, и враћа још увек празну мапу.

Да бисмо то поправили, можемо само сачекати док сви гороутини ураде свој посао, а затим повратак. Требало би две секунде, зар не?

```go
package concurrency

import "time"

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()
	}

	time.Sleep(2 * time.Second)

	return results
}
```

Сада када покренемо тестове које добијате (или не добијате - погледајте горе):

```sh
--- FAIL: TestCheckWebsites (0.00s)
        CheckWebsites_test.go:31: Wanted map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], got map[waat://furhurterwe.geds:false]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s
```

Ово није сјајно - зашто само један резултат? Можда бисмо покушали да поправимо ово повећањем време које чекамо - пробајте ако желите. Неће успети. Овде је проблем у томе променљива `url` се поново користи за сваку итерацију `for` петље - потребно је сваки пут нова вредност из `urls`. Али свака наша гороутина има референцу на променљиву `url` - немају своју независну копију. Па јесу _све_ писање вредности коју `for` има на крају итерације - последње урл. Због тога је једини резултат који имамо последњи урл.

Да бисте то поправили:

```go
package concurrency

import (
	"time"
)

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func(u string) {
			results[u] = wc(u)
		}(url)
	}

	time.Sleep(2 * time.Second)

	return results
}
```

By giving each anonymous function a parameter for the url - `u` - and then
calling the anonymous function with the `url` as the argument, we make sure that
the value of `u` is fixed as the value of `url` for the iteration of the loop
that we're launching the goroutine in. `u` is a copy of the value of `url`, and
so can't be changed.

Now if you're lucky you'll get:


Давањем сваке анонимне функције параметар за урл - `u` - и затим
позивајући анонимну функцију са `url` као аргументом, ми то осигуравамо
вредност `u` је фиксирана као вредност `url` за итерацију петље
да лансирамо гороутину у. `u` је копија вредности `url`, и
па се не може променити.

Ако будете имали среће, добићете:

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        2.012s
```

Али ако немате среће (ово је вероватније ако их покренете са референтном вредношћу јер ћете добити више покушаја)

```sh
fatal error: concurrent map writes

goroutine 8 [running]:
runtime.throw(0x12c5895, 0x15)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/panic.go:605 +0x95 fp=0xc420037700 sp=0xc4200376e0 pc=0x102d395
runtime.mapassign_faststr(0x1271d80, 0xc42007acf0, 0x12c6634, 0x17, 0x0)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:783 +0x4f5 fp=0xc420037780 sp=0xc420037700 pc=0x100eb65
github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1(0xc42007acf0, 0x12d3938, 0x12c6634, 0x17)
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x71 fp=0xc4200377c0 sp=0xc420037780 pc=0x12308f1
runtime.goexit()
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/asm_amd64.s:2337 +0x1 fp=0xc4200377c8 sp=0xc4200377c0 pc=0x105cf01
created by github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xa1

        ... many more scary lines of text ...
```

Ово је дуго и застрашујуће, али све што требамо је удахнути и прочитати стактрејс: `fatal error: concurrent map writes`. Понекад, када покренемо свој тестова, две гороутине уписују на мапу резултата тачно у исто време.
Мапе у Го-у не воле када им покушава да пише више ствари једном, и тако `fatal error`.

Ово је _раце цондитион_, грешка која се јавља када је излаз нашег софтвера зависно од времена и редоследа догађаја над којима немамо контролу.
Јер не можемо тачно контролисати када сваки гороутин уписује на мапу резултата, осетљиви смо на две гороутине које јој истовремено пишу.

Го нам може помоћи да уочимо услове трке помоћу уграђеног [_детектор расе_][godoc_race_detector].
Да бисте омогућили ову функцију, покрените тестове са заставицом `race` flag: `go test -race`.

Требали бисте добити излаз који изгледа овако:

```sh
==================
WARNING: DATA RACE
Write at 0x00c420084d20 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Previous write at 0x00c420084d20 by goroutine 7:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Goroutine 8 (running) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c

Goroutine 7 (finished) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c
==================
```

Појединости је, опет, тешко прочитати - али `WARNING: DATA RACE` је прилично недвосмислено. Читајући у тело грешке можемо видети две различите гороутине које изводе записе на мапи:

`Write at 0x00c420084d20 by goroutine 8:`

пише у исти блок меморије као и

`Previous write at 0x00c420084d20 by goroutine 7:`

Поврх тога, можемо видети линију кода у коју се уписује:

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12`

и линија кода у којој су покренути програми 7 и 8:

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11`

Све што требате знати штампа се на вашем терминалу - све што морате учинити је будите довољно стрпљиви да то прочитате.

### Канали

Ову трку података можемо решити координацијом наших програма користећи _канале_. Канали су Го дата структура која може и примати и слати вредности. Ове операције, заједно са њиховим детаљима, омогућавају комуникацију између различитих процеси.

У овом случају желимо да размислимо о комуникацији између родитељског процеса и сваки од програма који чини да изврши посао вођења Функција `WebsiteChecker` са УРЛ-ом.

```go
package concurrency

type WebsiteChecker func(string) bool
type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {
		go func(u string) {
			resultChannel <- result{u, wc(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		r := <-resultChannel
		results[r.string] = r.bool
	}

	return results
}
```

Поред мапе `results` сада имамо и `resultChannel` у који `make` на исти начин. `chan result` је врста канала - канал` result`.
Нови тип, `резултат` је направљен да повеже повратну вредност `WebsiteChecker` са УРЛ-ом који се проверава - то је структура од` string` и `bool`. Како нам није потребна ниједна вредност за именовање, свака од њих је анонимна унутар структуре; ово може бити корисно када је тешко знати како именовати вредност.

Сада када прелазимо преко УРЛ-ова, уместо да директно пишемо на `map` шаљемо структуру `result` за сваки позив` wc` на `resultChannel` са _сенд изјавом_. Ово користи оператора `<-`, узимајући канал на лево и вредност десно:

```go
// Send statement
resultChannel <- result{u, wc(u)}
```

Следећа петља `for` понавља се једном за сваки УРЛ. Унутра користимо _примити израз_, који додељује вредност примљену од канала променљива. Ово такође користи оператор `<-`, али сада са два операнда обрнуто: канал је сада на десној страни и променљива која додељујемо на левој страни:

```go
// Receive expression
r := <-resultChannel
```

Затим користимо примљени `result` за ажурирање мапе.

Слањем резултата у канал можемо да контролишемо време сваког писања у мапу резултата, осигуравајући да се то дешава једно по једно. Иако је сваки од позиви `wc` и свако слање на канал резултата одвијају се паралелно у оквиру свог процеса, сваки од резултата се обрађује један по један као вадимо вредности из канала резултата са изразом за пријем.

Паралелно смо упоредили део кода који смо желели да направимо брже пазећи да се део који се не може паралелно дешавати и даље дешава линеарно.
И комуницирали смо кроз вишеструке процесе који су укључени коришћењем канали.

Када покренемо референтну вредност:

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v2
BenchmarkCheckWebsites-8             100          23406615 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        2.377s
```
23406615 наносекунде - 0,023 секунде, отприлике сто пута брже од оригинална функција. Велики успех.

## Окончање

Ова вежба је била мало лакша на ТДД-у него обично. На неки начин јесмо учествовао у једном дугом рефакторирању функције `CheckWebsites`; улази и излази се никада нису мењали, већ су постали бржи. Али тестови које смо имали место, као и референтна тачка коју смо написали, омогућили су нам да рефакторизујемо `CheckWebsites` на начин који је задржао уверење да софтвер и даље ради, док демонстрирајући да је заправо постало брже.

Убрзавајући то, учили смо о томе

- *гороутине*, основна јединица подударности у Го-у, која нам омогућава да проверимо више више од једне веб странице истовремено.
- *анонимне функције*, које смо користили за покретање сваког од истовремених процеса који проверавају веб локације.
- *канали*, који помажу у организацији и контроли комуникације између различити процеси, што нам омогућава да избегнемо *грешку расе*.
- *детектор трке* који нам је помогао да решимо проблеме са истовременим кодом

### Чине га брзо

Једна формулација агилног начина израде софтвера, често погрешно приписана Кенту Беку, је:

> [Нека то успе, поправи, учини брзо][врф]

Тамо где „рад“ чини да тестови прођу, 'right' је рефакторисање кода и 'fast' је оптимизација кода како би се, на пример, брзо покренуо. Можемо само „убрзајте“ након што смо то успели и исправили. Имали смо среће да је
код који смо добили већ је доказано да ради и није требало реконструисан. Никада не бисмо требали покушавати да „убрзамо“ пре друга два корака су изведени јер

> [Превремена оптимизација је корен свега зла] [попт]
> - Доналд Кнутх

[DI]: dependency-injection.md
[врф]: http://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast
[godoc_детектор_расе]: https://blog.golang.org/race-detector
[попт]: http://wiki.c2.com/?PrematureOptimization
