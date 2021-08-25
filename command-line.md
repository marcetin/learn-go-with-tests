# Командна линија и структура пројекта

**[Сав код за ово поглавље можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/command-line)**

Наш власник производа сада жели _пивот_ увођењем друге апликације - апликације командне линије.

За сада ће само бити потребно да се забележи играчева победа када корисник откуца `Ruth wins`. Намера је да на крају постане алат за помоћ корисницима у игрању покера.

Власник производа жели да се база података дели између две апликације, тако да се лига ажурира према победама забележеним у новој апликацији.

## Подсетник на код

Имамо апликацију са `main.go` датотеком која покреће ХТТП сервер. ХТТП сервер нам неће бити занимљив за ову вежбу, али апстракција коју користи хоће. Зависи од `PlayerStore`-а.

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}
```

У претходном поглављу направили смо `FileSystemPlayerStore` који имплементира тај интерфејс. Требали бисмо бити у могућности поново користити ово за нашу нову апликацију.

## Прво преобликовање пројеката

Наш пројекат сада треба да створи две бинарне датотеке, наш постојећи веб сервер и апликацију за командну линију.

Пре него што се заглавимо у новом послу, требало би да структурирамо наш пројекат да то прилагоди.

До сада је сав код живио у једној фасцикли, на путањи која изгледа овако

`$GOPATH/src/github.com/your-name/my-app`

Да бисте могли да направите апликацију у Го -у, потребна вам је функција `main` унутар `package main`. До сада је сав наш "домен" код живео у `package main` и наш `func main` се може позивати на све.

Ово је до сада било у реду и добра је пракса да не претерујете са структуром пакета. Ако одвојите мало времена да прегледате стандардну библиотеку, видећете врло мало у смислу пуно фасцикли и структуре.

Срећом, прилично је једноставно додати структуру _ када вам затреба_.

Унутар постојећег пројекта креирајте `func main` директоријум са директоријем` webserver` унутар њега (нпр. `mkdir -p cmd/webserver`).

Померите `main.go` унутра.

Ако имате `tree` инсталирано, требали бисте га покренути и ваша структура би требала изгледати овако

```
.
├── file_system_store.go
├── file_system_store_test.go
├── cmd
│   └── webserver
│       └── main.go
├── league.go
├── server.go
├── server_integration_test.go
├── server_test.go
├── tape.go
└── tape_test.go
```

Сада заправо имамо раздвајање између наше апликације и библиотечког кода, али сада морамо да променимо нека имена пакета. Запамтите да када правите Го апликацију њен пакет _мора_ бити `main`.

Промените све остале кодове да бисте добили пакет под називом `poker`.

Коначно, морамо да увозимо овај пакет у `main.go` тако да га можемо користити за креирање нашег веб сервера. Тада можемо користити наш библиотечки код помоћу `poker.FunctionName`.

Путања ће бити различита на вашем рачунару, али би требало да буде слична овој:

```go
//cmd/webserver/main.go
package main

import (
	"github.com/marcetin/nauci-go-sa-testovima/command-line/v1"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := poker.NewPlayerServer(store)

    log.Fatal(http.ListenAndServe(":5000", server))
}
```

Читав пут може изгледати помало узнемирујуће, али овако можете увести _било коју_ јавно доступну библиотеку у свој код.

Одвајањем нашег кода домена у посебан пакет и предавањем на јавни репо, попут ГитХуб -а, сваки Го програмер може написати свој код који увози тај пакет доступних функција. Први пут када покушате да га покренете жалиће се да не постоји, али све што треба да урадите је да покренете `go get`.

Осим тога, корисници могу прегледати [докумантацију на godoc.org](https://godoc.org/github.com/marcetin/nauci-go-sa-testovima/command-line/v1).

### Завршне провере

- Унутар роот-а покрените `go test` и проверите да ли још увек пролази
- Уђите у наш `cmd/webserver` и `go run main.go`
    - Посетите `http://localhost:5000/league` и требало би да видите да и даље ради

### Пешачки костур

Пре него што се заглавимо у писању тестова, додајмо нову апликацију коју ће наш пројекат изградити. Направите још један директоријум унутар `cmd` -а под називом` cli` (интерфејс командне линије) и додајте `main.go` са следећим

```go
//cmd/cli/main.go
package main

import "fmt"

func main() {
	fmt.Println("Let's play poker")
}
```

Први услов који ћемо решити је бележење победе када корисник откуца `{PlayerName} wins`.

## Прво напишите тест

Знамо да морамо да направимо нешто што се зове `CLI` што ће нам омогућити да` Play` покер. Мораће да прочита унос корисника, а затим забележи победе у `PlayerStore`.

Пре него што одемо предалеко, само да напишемо тест да проверимо да ли се интегрише са `PlayerStore` -ом онако како бисмо желели.

Унутар `CLI_test.go` (у корену пројекта, не унутар `cmd`)

```go
//CLI_test.go
package poker

import "testing"

func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	cli := &CLI{playerStore}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}
}
```

- Можемо користити наш `StubPlayerStore` из других тестова
- Своју зависност преносимо у наш још не постојећи тип `CLI`
- Покрените игру по неписаној `PlayPoker` методи
- Проверите да ли је забележена победа

## Покушајте да покренете тест

```
# github.com/marcetin/nauci-go-sa-testovima/command-line/v2
./cli_test.go:25:10: undefined: CLI
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

У овом тренутку би вам требало бити довољно удобно да креирате нашу нову `CLI` структуру са одговарајућим пољем за нашу зависност и додате методу.

Требало би да завршите са оваквим кодом

```go
//CLI.go
package poker

type CLI struct {
	playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {}
```

Упамтите да само покушавамо покренути тест како бисмо могли провјерити да тест није успио како смо се надали

```
--- FAIL: TestCLI (0.00s)
    cli_test.go:30: expected a win call but didn't get any
FAIL
```

## Напишите довољно кода да прође

```go
//CLI.go
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
```

То би требало да прође.

Затим морамо симулирати читање са `Stdin` (унос од корисника) како бисмо могли забиљежити побједе за одређене играче.

Хајде да проширимо наш тест на ово.

## Прво напишите тест

```go
//CLI_test.go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}

	got := playerStore.winCalls[0]
	want := "Chris"

	if got != want {
		t.Errorf("didn't record correct winner, got %q, want %q", got, want)
	}
}
```

`os.Stdin` је оно што ћемо користити у `main` за хватање уноса корисника. То је `*File` испод хаубе, што значи да имплементира `io.Reader` који је до сада познат згодан начин за снимање текста.

Ми стварамо `io.Reader` у нашем тесту помоћу практичних `strings.NewReader`, испуњавајући га оним што очекујемо од корисника да откуца.

## Покушајте да покренете тест

`./CLI_test.go:12:32: too many values in struct initializer`

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

Морамо додати нашу нову зависност у `CLI`.

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}
```

## Напишите довољно кода да прође

```
--- FAIL: TestCLI (0.00s)
    CLI_test.go:23: didn't record the correct winner, got 'Cleo', want 'Chris'
FAIL
```

Не заборавите да прво учините најједноставнију ствар

```go
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Chris")
}
```

Тест пролази. Додаћемо још један тест који ће нас натерати да следећи напишемо прави код, али прво, хајде да преуредимо.

## Рефактор

У `server_test` смо раније проверили да ли су победе забележене као што имамо овде. СУШИМО ту тврдњу у помоћника

```go
//server_test.go
func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}
```

Сада замените тврдње у `server_test.go` и `CLI_test.go`.

Тест би сада требао да изгледа овако

```go
//CLI_test.go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	assertPlayerWin(t, playerStore, "Chris")
}
```

Хајде сада да напишемо још један тест са различитим корисничким уносом који ће нас натерати да га заиста читамо.

## Прво напишите тест

```go
//CLI_test.go
func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Cleo")
	})

}
```

## Покушајте да покренете тест

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/record_chris_win_from_user_input
    --- PASS: TestCLI/record_chris_win_from_user_input (0.00s)
=== RUN   TestCLI/record_cleo_win_from_user_input
    --- FAIL: TestCLI/record_cleo_win_from_user_input (0.00s)
        CLI_test.go:27: did not store correct winner got 'Chris' want 'Cleo'
FAIL
```

## Напишите довољно кода да прође

Користићемо [`bufio.Scanner`](https://golang.org/pkg/bufio/) за читање уноса из `io.Reader`.

> Бufio пакет имплементира међуспремник I/O. Обухвата io.Reader или io.Writer објекат, стварајући други објекат (Reader или Writer) који такође имплементира интерфејс, али пружа баферисање и одређену помоћ за текстуални `I/O`.

Ажурирајте код на следеће

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
```

Сада ће тестови проћи.

- `Scanner.Scan()` ће читати до новог реда.
- Затим користимо `Scanner.Text()` да вратимо `string` у који је скенер читао.

Сада када имамо неке положене тестове, требали бисмо ово повезати у `main`. Упамтите да бисмо увек требали настојати да имамо што интегрисанији радни софтвер што је брже могуће.

У `main.go` додајте следеће и покрените га. (можда ћете морати да прилагодите путању друге зависности тако да одговара ономе што је на вашем рачунару)

```go
package main

import (
	"fmt"
	"github.com/marcetin/nauci-go-sa-testovima/command-line/v3"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	game := poker.CLI{store, os.Stdin}
	game.PlayPoker()
}
```

Требало би да добијете грешку

```
command-line/v3/cmd/cli/main.go:32:25: implicit assignment of unexported field 'playerStore' in poker.CLI literal
command-line/v3/cmd/cli/main.go:32:34: implicit assignment of unexported field 'in' in poker.CLI literal
```

Оно што се овде дешава је зато што покушавамо да доделимо поља `playerStore` и` in` у `CLI`. То су неекспортована (приватна) поља. Ово можемо _учинити_ у нашем тестном коду јер је наш тест у истом пакету као `CLI` (` poker`). Али наш `main` је у пакету` main` па нема приступ.

Ово наглашава важност _интегрисања вашег рада_. С правом смо учинили зависности нашег `CLI`-а приватним (јер не желимо да буду изложени корисницима` "CLI"-а), али нисмо направили начин да га корисници конструишу.

Постоји ли начин да се овај проблем ухвати раније?

### `package mypackage_test`

У свим осталим примерима до сада, када правимо пробну датотеку, изјављујемо да се налази у истом пакету који тестирамо.

Ово је у реду и значи да у чудним приликама када желимо да тестирамо нешто интерно у пакету имамо приступ неекспортованим типовима.

Али с обзиром да смо се залагали за _не_ тестирање унутрашњих ствари _опћенито_, може ли Го помоћи да се то спроведе? Шта ако бисмо могли да тестирамо наш код где имамо приступ само извезеним типовима (попут нашег `main`-а)?

Када пишете пројекат са више пакета, топло бих вам препоручио да назив вашег тестног пакета има `_test` на крају. Када то учините, моћи ћете да имате приступ само јавним типовима у свом пакету. Ово би помогло у овом конкретном случају, али такође помаже у примени дисциплине само тестирања јавних АПИ -ја. Ако и даље желите да тестирате унутрашњост, можете направити посебан тест са пакетом који желите да тестирате.

Изрека за ТДД је да ако не можете да тестирате свој код, онда ће корисницима вашег кода вероватно бити тешко да се интегришу са њим. Употреба `package foo_test` ће вам помоћи у томе, приморавши вас да тестирате свој код као да га увозите као што ће то учинити корисници вашег пакета.

Пре него што поправимо `маин`, променимо пакет нашег теста унутар `CLI_test.go` у `poker_test`.

Ако имате добро конфигурисан ИДЕ, одједном ћете видети много црвене боје! Ако покренете компајлер, добићете следеће грешке

```
./CLI_test.go:12:19: undefined: StubPlayerStore
./CLI_test.go:17:3: undefined: assertPlayerWin
./CLI_test.go:22:19: undefined: StubPlayerStore
./CLI_test.go:27:3: undefined: assertPlayerWin
```

Сада смо наишли на још питања о дизајну паковања. Да бисмо тестирали наш софтвер, направили смо неекспортиране стубове и помоћне функције које нам више нису доступне за употребу у нашем `CLI_test` јер су помагачи дефинисани у `_test.go` датотекама у `poker` пакету.

#### Да ли желимо да наши стубови и помагачи буду „јавни“?

Ово је субјективна дискусија. Могло би се тврдити да не желите да загађујете АПИ свог пакета кодом ради олакшавања тестова.

У презентацији ["Напредно тестирање уз Го"](https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=53) аутора Митчела Хашимотоа, описано је како се у ХасхиЦорп-у залажу за то како би корисници пакета могли писати тестове без потребе за поновним проналажењем точкића за писање котача. У нашем случају, ово би значило да било ко ко користи наш `poker` пакет неће морати да креира сопствени стуб `PlayerStore` ако жели да ради са нашим кодом.

Анегдотски сам користио ову технику у другим дељеним пакетима и показала се изузетно корисном у смислу уштеде времена корисника приликом интеграције са нашим пакетима.

Хајде да направимо датотеку под називом `testing.go` и додамо наш стуб и наше помоћнике.

```go
//testing.go
package poker

import "testing"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}

// todo for you - the rest of the helpers
```

Морате да учините помоћнике јавним (запамтите да се извоз врши великим словом на почетку) ако желите да буду изложени увозницима нашег пакета.

У нашем `CLI` тесту морате да позовете код као да га користите у другом пакету.

```go
//CLI_test.go
func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.CLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.CLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

}
```

Сада ћете видети да имамо исте проблеме као и у `main`-у

```
./CLI_test.go:15:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:15:39: implicit assignment of unexported field 'in' in poker.CLI literal
./CLI_test.go:25:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:25:39: implicit assignment of unexported field 'in' in poker.CLI literal
```

Најлакши начин да то заобиђете је да направите конструктор какав имамо за друге типове. Такође ћемо променити `CLI` тако да он складишти `bufio.Scanner` уместо читача јер је сада аутоматски умотан у време изградње.

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}
```

Радећи ово, тада можемо поједноставити и прерадити наш код за читање

```go
//CLI.go
func (cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
```

Промените тест да бисте уместо њега користили конструктор и требало би да се вратимо на пролаз тестова.

Коначно, можемо се вратити на наш нови `main.go` и користити конструктор који смо управо направили

```go
//cmd/cli/main.go
game := poker.NewCLI(store, os.Stdin)
```

Покушајте да га покренете, откуцајте "Bob wins".

### Рефактор

Имамо неколико понављања у нашим одговарајућим апликацијама где отварамо датотеку и правимо `file_system_store` од њеног садржаја. Ово се осећа као мала слабост у дизајну нашег пакета, па бисмо у њему требали направити функцију која инкапсулира отварање датотеке са путање и враћање `PlayerStore`-а.

```go
//file_system_store.go
func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, closeFunc, nil
}
```

Сада прерадите обе наше апликације да користе ову функцију за креирање продавнице.

#### CLI код апликације

```go
//cmd/cli/main.go
package main

import (
	"fmt"
	"github.com/marcetin/nauci-go-sa-testovima/command-line/v3"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	poker.NewCLI(store, os.Stdin).PlayPoker()
}
```

#### Код апликације веб сервера

```go
//cmd/webserver/main.go
package main

import (
	"github.com/marcetin/nauci-go-sa-testovima/command-line/v3"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

Уочите симетрију: упркос различитим корисничким интерфејсима, подешавање је готово идентично. Ово се чини као добра потврда нашег дизајна до сада.
Такође приметите да `FileSystemPlayerStoreFromFile` враћа функцију затварања, тако да можемо затворити основну датотеку када завршимо са коришћењем продавнице.

## Окончање

### Структура пакета

Ово поглавље је значило да желимо да направимо две апликације, користећи поново код домена који смо до сада написали. Да бисмо то учинили, морали смо да ажурирамо нашу структуру пакета тако да имамо посебне фасцикле за наш `main`.

Тиме смо наишли на проблеме интеграције због неекспортованих вредности, па ово додатно показује вредност рада у малим „кришкама“ и често интегрисање.

Научили смо како нам `mypackage_test` помаже у стварању окружења за тестирање које је исто искуство за остале пакете који се интегришу са вашим кодом, како би вам помогло да ухватите проблеме при интеграцији и видите са колико је лако (или не!) Радити са вашим кодом.

### Читање корисничког уноса

Видели смо како нам је читање из `os.Stdin` -а веома лако радити јер имплементира `io.Reader`. Користили смо `bufio.Scanner` за лако читање уноса корисника по редак.

### Једноставне апстракције доводе до једноставније поновне употребе кода

Готово да није било напора да интегришемо `PlayerStore` у нашу нову апликацију (након што смо извршили прилагођавања пакета), а касније је и тестирање било врло једноставно јер смо одлучили да откријемо и нашу верзију стуб -а.
