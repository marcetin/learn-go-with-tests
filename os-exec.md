# OS Exec

**[Сав код можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/q-and-a/os-exec)**

[keith6014](https://www.reddit.com/user/keith6014) пита на [reddit](https://www.reddit.com/r/golang/comments/aaz8ji/testdata_and_function_setup_help/)

> Извршавам команду помоћу os/exec.Command() која је генерисала "XML" податке. Команда ће се извршавати у функцији која се зове GetData().

> Да бих тестирао GetData(), имам неке тестне податке које сам направио.

> У мом _test.go имам ТестГетДата који позива GetData(), али који ће користити os.exec, уместо тога бих волео да он користи моје тестне податке.

> Који је добар начин да се то постигне? Приликом позивања GetData требам ли имати "тест" начин означавања тако да чита датотеку, тј. GetData (низ начина)?

Неколико ствари

- Када је нешто тешко тестирати, то је често због тога што раздвајање брига није сасвим у реду
- Не додавајте "тест моде" у свој код, уместо тога користите [Dependency Injection](/dependency-injection.md) тако да можете моделирати своје зависности и одвојити бриге.

Узео сам слободу да погодим како би код могао изгледати


```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData() string {
	cmd := exec.Command("cat", "msg.xml")

	out, _ := cmd.StdoutPipe()
	var payload Payload
	decoder := xml.NewDecoder(out)

	// these 3 can return errors but I'm ignoring for brevity
	cmd.Start()
	decoder.Decode(&payload)
	cmd.Wait()

	return strings.ToUpper(payload.Message)
}
```

- Користи `exec.Command` који вам омогућава извршавање спољне команде процеса
- Снимамо излаз у `cmd.StdoutPipe` који нам враћа `io.ReadCloser` (ово ће постати важно)
- Остатак кода је мање -више копиран и залепљен из [одличне документације](https://golang.org/pkg/os/exec/#example_Cmd_StdoutPipe).
    - Снимимо било који излаз са стдоут -а у `io.ReadCloser` и затим `Start` команду, а затим сачекамо да се сви подаци прочитају позивањем `Wait`. Између та два позива декодирамо податке у нашу `Payload` структуру.

Here is what is contained inside `msg.xml`

```xml
<payload>
    <message>Happy New Year!</message>
</payload>
```

Написао сам једноставан тест да то покажем на делу

```go
func TestGetData(t *testing.T) {
	got := GetData()
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

## Код који се може тестирати

Код за тестирање је одвојен и има једну намену. Мени се чини да постоје два главна проблема за овај код

1. Преузимање сирових "XML" података
2. Декодирање "XML" података и примена наше пословне логике (у овом случају `strings.ToUpper` на` <message> `)

Први део је само копирање примера из стандардног либ -а.

Други део је где имамо своју пословну логику и гледајући код можемо видети где почиње "шав" у нашој логици; ту добијамо `io.ReadCloser`. Можемо користити ову постојећу апстракцију да одвојимо забринутости и учинимо наш код тестираним.

**Проблем са ГетДата -ом је што је пословна логика повезана са средствима за добијање "XML"-а. Да бисмо наш дизајн учинили бољим, морамо их одвојити**

Наши `TestGetData` могу деловати као наш тест интеграције између наше две бриге, па ћемо се тога држати како бисмо били сигурни да наставља да ради.

Ево како изгледа ново раздвојени код

```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := ioutil.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

func TestGetDataIntegration(t *testing.T) {
	got := GetData(getXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

Сада када `GetData` узима своје уносе само из `io.Reader`-а, учинили смо га тестираним и више се не брине како се подаци преузимају; људи могу поново користити функцију са било чим што враћа `io.Reader` (што је изузетно уобичајено). На пример, могли бисмо да почнемо са преузимањем "XML" -а са УРЛ -а уместо из командне линије.

```go
func TestGetData(t *testing.T) {
	input := strings.NewReader(`
<payload>
    <message>Cats are the best animal</message>
</payload>`)

	got := GetData(input)
	want := "CATS ARE THE BEST ANIMAL"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

```

Ево примера јединичног теста за `GetData`.

Одвајањем брига и употребом постојећих апстракција у оквиру Го тестирања, наша важна пословна логика је лака.
