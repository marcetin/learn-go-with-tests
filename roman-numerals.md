# Roman Numerals

**[Сав код за ово поглавље можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/roman-numerals)**

Написаћемо функцију која претвара [Римску каталошку нумеру](http://codingdojo.org/kata/RomanNumerals/) као део процеса интервјуа. Ово поглавље ће показати како се с њим можете носити са ТДД-ом.

We are going to write a function which converts an [Арапски број](https://en.wikipedia.org/wiki/Arabic_numerals) (бројеве од 0 до 9) у римски број.

Ако нисте чули за [Римске бројеве](https://en.wikipedia.org/wiki/Roman_numerals), Римљани су записали бројеве.

Градите их лепљењем симбола и ти симболи представљају бројеве

Дакле, `I` је" једно ". `III` је три.

Изгледа лако, али постоји неколико занимљивих правила. `V` значи пет, али `IV` је 4 (не `IIII`).

`MCMLXXXIV` је 1984. То изгледа компликовано и тешко је замислити како можемо написати код да то схватимо од самог почетка.

Као што наглашава ова књига, кључна вештина програмера је да покушају да идентификују „танке вертикалне кришке“ _корисне_ функционалности, а затим ** понављају **. ТДД радни ток помаже у олакшавању итеративног развоја.

Дакле, радије него 1984, почнимо са 1.

## Прво напишите тест

```go
func TestRomanNumerals(t *testing.T) {
	got := ConvertToRoman(1)
	want := "I"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

Ако сте у књизи стигли толико далеко, надам се да вам се чини врло досадним и рутинским. То је добра ствар.

## Покушајте да покренете тест

`./numeral_test.go:6:9: undefined: ConvertToRoman`

Нека преводилац води пут

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

Креирајте нашу функцију, али још увек немојте да прођете тест, увек проверите да тестови не успеју онако како очекујете

```go
func ConvertToRoman(arabic int) string {
	return ""
}
```

Требало би да ради сада

```go
=== RUN   TestRomanNumerals
--- FAIL: TestRomanNumerals (0.00s)
    numeral_test.go:10: got '', want 'I'
FAIL
```

## Напишите довољно кода да прође

```go
func ConvertToRoman(arabic int) string {
	return "I"
}
```

## Рефактор

Још увек није много за рефакторирање.

_Знам_ чудно ми је само кодирање резултата, али са ТДД-ом желимо да будемо ван „црвеног“ што је дуже могуће. Можда се осећа _као да нисмо постигли много, али дефинисали смо наш АПИ и добили тест који обухвата једно од наших правила; чак и ако је „прави“ код прилично глуп.

Сада искористите тај нелагодан осећај да напишете нови тест који ће нас присилити да напишемо мало мање глупи код.

## Прво напишите тест

Можемо да користимо подтестове за лепо груписање тестова

```go
func TestRomanNumerals(t *testing.T) {
	t.Run("1 gets converted to I", func(t *testing.T) {
		got := ConvertToRoman(1)
		want := "I"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("2 gets converted to II", func(t *testing.T) {
		got := ConvertToRoman(2)
		want := "II"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
```

## Покушајте да покренете тест

```
=== RUN   TestRomanNumerals/2_gets_converted_to_II
    --- FAIL: TestRomanNumerals/2_gets_converted_to_II (0.00s)
        numeral_test.go:20: got 'I', want 'II'
```

Нема пуно изненађења тамо

## Напишите довољно кода да прође

```go
func ConvertToRoman(arabic int) string {
	if arabic == 2 {
		return "II"
	}
	return "I"
}
```

Да, и даље се чини да се заправо не бавимо проблемом. Зато морамо да напишемо још тестова који ће нас водити напред.

## Рефактор

Имамо неколико понављања у нашим тестовима. Када тестирате нешто што се осећа као да је реч о „датом уносу X, очекујемо Y“, вероватно бисте требали да користите тестове засноване на табелама.

```go
func TestRomanNumerals(t *testing.T) {
	cases := []struct {
		Description string
		Arabic      int
		Want        string
	}{
		{"1 gets converted to I", 1, "I"},
		{"2 gets converted to II", 2, "II"},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got := ConvertToRoman(test.Arabic)
			if got != test.Want {
				t.Errorf("got %q, want %q", got, test.Want)
			}
		})
	}
}
```

Сада можемо лако додати још случајева без потребе за писањем било ког пробног узорка.

Идемо даље и идемо на 3

## Прво напишите тест

Нашим случајевима додајте следеће

```go
{"3 gets converted to III", 3, "III"},
```

## Покушајте да покренете тест

```
=== RUN   TestRomanNumerals/3_gets_converted_to_III
    --- FAIL: TestRomanNumerals/3_gets_converted_to_III (0.00s)
        numeral_test.go:20: got 'I', want 'III'
```

## Напишите довољно кода да прође

```go
func ConvertToRoman(arabic int) string {
	if arabic == 3 {
		return "III"
	}
	if arabic == 2 {
		return "II"
	}
	return "I"
}
```

## Рефактор

ОК, почињем да не уживам у овим иф изјавама и ако довољно добро погледате код, видећете да градимо низ `I` на основу величине „арапског“.

„Знамо“ да ћемо за сложеније бројеве радити неку врсту аритметике и спајања низова.

Покушајмо са рефактором имајући на уму ове мисли, он _можда не би_ био погодан за крајње решење, али то је у реду. Увек можемо да бацимо свој код и започнемо изнова са тестовима којима се морамо водити.

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for i:=0; i<arabic; i++ {
		result.WriteString("I")
	}

	return result.String()
}
```

Можда нисте користили [`strings.Builder`](https://golang.org/pkg/strings/#Builder) раније.

> "Builder" се користи за ефикасну изградњу низа помоћу метода Врите. Минимизира копирање меморије.

Обично се не бих замарао таквом оптимизацијом док не будем имао стварни проблем са перформансама, али количина кода није много већа од „ручног“ додавања на стринг, па бисмо могли користити и бржи приступ.

Код ми изгледа боље и описује домен _као што га тренутно знамо_.

### Римљани су такође били на СУВОМ ...

Ствари сада почињу да се компликују. Римљани су у својој мудрости мислили да ће понављање ликова постати тешко читати и бројати. Дакле, правило римских бројева је да не можете да поновите исти лик више од 3 пута заредом.

Уместо тога, узимате следећи највиши симбол, а затим "одузимате" стављањем симбола лево од њега. Не могу се сви симболи користити као одузимачи; само I (1), X (10) and C (100).


На пример, `5` у римским бројевима је `V`. Да бисте креирали 4, не радите `IIII`, већ `IV`.

## Прво напишите тест

```go
{"4 gets converted to IV (can't repeat more than 3 times)", 4, "IV"},
```

## Покушајте да покренете тест

```
=== RUN   TestRomanNumerals/4_gets_converted_to_IV_(cant_repeat_more_than_3_times)
    --- FAIL: TestRomanNumerals/4_gets_converted_to_IV_(cant_repeat_more_than_3_times) (0.00s)
        numeral_test.go:24: got 'IIII', want 'IV'
```

## Напишите довољно кода да прође

```go
func ConvertToRoman(arabic int) string {

	if arabic == 4 {
		return "IV"
	}

	var result strings.Builder

	for i:=0; i<arabic; i++ {
		result.WriteString("I")
	}

	return result.String()
}
```

## Рефактор

Не "свиђа ми се" што смо прекинули наш образац грађења низова и желим да наставим са тим.

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for i := arabic; i > 0; i-- {
		if i == 4 {
			result.WriteString("IV")
			break
		}
		result.WriteString("I")
	}

	return result.String()
}
```

Да би се 4 "уклопило" са мојим тренутним размишљањем, сада одбројавам од арапског броја, додајући симболе у наш низ како напредујемо. Нисам сигуран да ли ће ово дугорочно успети, али да видимо!

Направимо 5 послова

## Прво напишите тест

```go
{"5 gets converted to V", 5, "V"},
```

## Покушајте да покренете тест

```
=== RUN   TestRomanNumerals/5_gets_converted_to_V
    --- FAIL: TestRomanNumerals/5_gets_converted_to_V (0.00s)
        numeral_test.go:25: got 'IIV', want 'V'
```

## Напишите довољно кода да прође

Само копирајте приступ који смо урадили за 4

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for i := arabic; i > 0; i-- {
		if i == 5 {
			result.WriteString("V")
			break
		}
		if i == 4 {
			result.WriteString("IV")
			break
		}
		result.WriteString("I")
	}

	return result.String()
}
```

## Рефактор

Понављање у оваквим петљама обично је знак апстракције која чека на прозивање. Петље кратког споја могу бити ефикасан алат за читљивост, али такође вам могу рећи нешто друго.

Ми петљамо преко свог арапског броја и ако притиснемо одређене симболе, називамо `break`, али оно што стварно _радимо_ је одузимање преко `i` на шункаст начин.

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for arabic > 0 {
		switch {
		case arabic > 4:
			result.WriteString("V")
			arabic -= 5
		case arabic > 3:
			result.WriteString("IV")
			arabic -= 4
		default:
			result.WriteString("I")
			arabic--
		}
	}

	return result.String()
}

```

- С обзиром на сигнале које читам из нашег кода, на основу тестова неких врло основних сценарија, видим да за израду римског броја морам да одузмем од `арапског` док примењујем симболе
- Петља `for` се више не ослања на `i` и уместо тога наставићемо да градимо свој низ све док од арапског не одузмемо довољно симбола.

Прилично сам сигуран да ће овај приступ важити и за 6 (VI), 7 (VII) и 8 (VIII). Без обзира на то додајте случајеве у наш тестни пакет и проверите (нећу укључити код за краткоћу, проверите гитхуб за узорке ако нисте сигурни).

9 следи исто правило као и 4 у томе што би требало да одузмемо `I` од приказа следећег броја. 10 је у римским бројевима представљено с `X`; па би зато 9 требало да буде `IX`.

## Прво напишите тест

```go
{"9 gets converted to IX", 9, "IX"}
```
## Покушајте да покренете тест

```
=== RUN   TestRomanNumerals/9_gets_converted_to_IX
    --- FAIL: TestRomanNumerals/9_gets_converted_to_IX (0.00s)
        numeral_test.go:29: got 'VIV', want 'IX'
```

## Напишите довољно кода да прође

Требали бисмо бити у могућности да усвојимо исти приступ као и раније

```go
case arabic > 8:
    result.WriteString("IX")
    arabic -= 9
```

## Рефактор

_Чини се_ као да нам код још увек говори да негде постоји рефактор, али то ми није потпуно очигледно, па наставимо даље.

Прескочићу и код за ово, али у своје тест случајеве додајте тест за `10` који треба да буде `X` и учините да прође пре него што прочитате даље.

Ево неколико тестова које сам додао јер сам уверен да би до 39 наш код требао радити

```go
{"10 gets converted to X", 10, "X"},
{"14 gets converted to XIV", 14, "XIV"},
{"18 gets converted to XVIII", 18, "XVIII"},
{"20 gets converted to XX", 20, "XX"},
{"39 gets converted to XXXIX", 39, "XXXIX"},
```

Ако сте икада радили ОО програмирање, знаћете да изјаве `switch` требате гледати с мало сумње. Обично хватате концепт или податке унутар неког императивног кода када би у ствари уместо тога могли да буду ухваћени у структуру класе.

Го није стриктно ОО, али то не значи да игноришемо лекције које ОО нуди у потпуности (онолико колико би неки желели да вам кажу).

Наша изјава о пребацивању описује неке истине о римским бројевима заједно са понашањем.

То можемо рефакторизирати раздвајањем података од понашања.

```go
type RomanNumeral struct {
	Value  int
	Symbol string
}

var allRomanNumerals = []RomanNumeral {
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for _, numeral := range allRomanNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}

	return result.String()
}
```

Осећам се много боље. Нека правила око бројева прогласили смо подацима, а не скривеним у алгоритму и можемо видети како само радимо кроз арапски број, покушавајући да додамо симболе нашем резултату ако одговарају.

Да ли ова апстракција делује за веће бројеве? Проширите тестни пакет тако да ради за римски број за 50 који је `L`.

Ево неколико тестова, покушајте и учините да прођу.

```go
{"40 gets converted to XL", 40, "XL"},
{"47 gets converted to XLVII", 47, "XLVII"},
{"49 gets converted to XLIX", 49, "XLIX"},
{"50 gets converted to L", 50, "L"},
```

Потребна помоћ? У овој суштини можете видети које симболе додати [овај гист](https://gist.github.com/pamelafox/6c7b948213ba55332d86efd0f0b037de).


## And the rest!

Ево преосталих симбола

| Арапски | Римски |
| ------- | :----: |
| 100     |    C   |
| 500     |    D   |
| 1000    |    M   |

Заузети исти приступ за преостале симболе, требало би само да се додају подаци и тестовима и нашем низу симбола.

Да ли ваш код ради за `1984`: `MCMLXXXIV`?

Ево моје последње пробне верзије

```go
func TestRomanNumerals(t *testing.T) {
	cases := []struct {
		Arabic int
		Roman  string
	}{
		{Arabic: 1, Roman: "I"},
		{Arabic: 2, Roman: "II"},
		{Arabic: 3, Roman: "III"},
		{Arabic: 4, Roman: "IV"},
		{Arabic: 5, Roman: "V"},
		{Arabic: 6, Roman: "VI"},
		{Arabic: 7, Roman: "VII"},
		{Arabic: 8, Roman: "VIII"},
		{Arabic: 9, Roman: "IX"},
		{Arabic: 10, Roman: "X"},
		{Arabic: 14, Roman: "XIV"},
		{Arabic: 18, Roman: "XVIII"},
		{Arabic: 20, Roman: "XX"},
		{Arabic: 39, Roman: "XXXIX"},
		{Arabic: 40, Roman: "XL"},
		{Arabic: 47, Roman: "XLVII"},
		{Arabic: 49, Roman: "XLIX"},
		{Arabic: 50, Roman: "L"},
		{Arabic: 100, Roman: "C"},
		{Arabic: 90, Roman: "XC"},
		{Arabic: 400, Roman: "CD"},
		{Arabic: 500, Roman: "D"},
		{Arabic: 900, Roman: "CM"},
		{Arabic: 1000, Roman: "M"},
		{Arabic: 1984, Roman: "MCMLXXXIV"},
		{Arabic: 3999, Roman: "MMMCMXCIX"},
		{Arabic: 2014, Roman: "MMXIV"},
		{Arabic: 1006, Roman: "MVI"},
		{Arabic: 798, Roman: "DCCXCVIII"},
	}
	for _, test := range cases {
		t.Run(fmt.Sprintf("%d gets converted to %q", test.Arabic, test.Roman), func(t *testing.T) {
			got := ConvertToRoman(test.Arabic)
			if got != test.Roman {
				t.Errorf("got %q, want %q", got, test.Roman)
			}
		})
	}
}
```

- Уклонио сам `description` јер сам осетио да _подаци_ описују довољно информација.
- Додао сам још неколико оштрих случајева које сам пронашао само да бих добио мало више самопоуздања. Са табеларним тестовима ово је врло јефтино урадити.

Нисам променио алгоритам, требало је само да ажурирам низ `allRomanNumerals`.

```go
var allRomanNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}
```

## Рашчлањивање римских бројева

Још нисмо завршили. Даље ћемо написати функцију која _из_ римског броја претвара у `int`


## Прво напишите тест

Овде можемо поново да користимо наше тест случајеве уз мало рефакторирања

Преместите променљиву `cases` изван теста као променљиву пакета у блок `var`.

```go
func TestConvertingToArabic(t *testing.T) {
	for _, test := range cases[:1] {
		t.Run(fmt.Sprintf("%q gets converted to %d", test.Roman, test.Arabic), func(t *testing.T) {
			got := ConvertToArabic(test.Roman)
			if got != test.Arabic {
				t.Errorf("got %d, want %d", got, test.Arabic)
			}
		})
	}
}
```

Приметите да за сада користим функцију пресека да бих само покренуо један од тестова (`cases[:1]`), јер је покушај да сви ти тестови прођу одједном превелики скок

## Покушајте да покренете тест

```
./numeral_test.go:60:11: undefined: ConvertToArabic
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

Додајте нашу нову дефиницију функције

```go
func ConvertToArabic(roman string) int {
	return 0
}
```

Тест би сада требао да се покрене и не успе

```
--- FAIL: TestConvertingToArabic (0.00s)
    --- FAIL: TestConvertingToArabic/'I'_gets_converted_to_1 (0.00s)
        numeral_test.go:62: got 0, want 1
```

## Напишите довољно кода да прође

Знате шта да радите

```go
func ConvertToArabic(roman string) int {
	return 1
}
```

Затим, промените индекс пресека у нашем тесту да бисте прешли на следећи тест случај (нпр. `cases[:2]`). Нека вам проследи најглупљи код који се можете сетити, наставите са писањем глупог кода (најбоља књига икад зар не?) И за трећи случај. Ево мог глупог кода.

```go
func ConvertToArabic(roman string) int {
	if roman == "III" {
		return 3
	}
	if roman == "II" {
		return 2
	}
	return 1
}
```

Кроз глупост _реалног кода који ради_ можемо почети да видимо образац као и пре. Морамо да прелистамо улаз и направимо _нешто_, у овом случају укупно.

```go
func ConvertToArabic(roman string) int {
	total := 0
	for range roman {
		total++
	}
	return total
}
```

## Прво напишите тест

Даље прелазимо на `cases[:4]` (`IV`) који сада не успева јер враћа 2 јер је то дужина низа.

## Напишите довољно кода да прође

```go
// earlier..
type RomanNumerals []RomanNumeral

func (r RomanNumerals) ValueOf(symbol string) int {
	for _, s := range r {
		if s.Symbol == symbol {
			return s.Value
		}
	}

	return 0
}

// later..
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		// look ahead to next symbol if we can and, the current symbol is base 10 (only valid subtractors)
		if i+1 < len(roman) && symbol == 'I' {
			nextSymbol := roman[i+1]

			// build the two character string
			potentialNumber := string([]byte{symbol, nextSymbol})

			// get the value of the two character string
			value := allRomanNumerals.ValueOf(potentialNumber)

			if value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++
			}
		} else {
			total++
		}
	}
	return total
}
```

Ово је ужасно, али успева. Толико је лоше да сам осетио потребу да додам коментаре.

- Желео сам да могу да пронађем целобројну вредност за дати римски број, па сам од нашег низа `RomanNumeral` направио тип, а затим му додао методу `ValueOf`
- Следеће у нашој петљи морамо гледати унапред _ако је низ довољно велик _а тренутни симбол је важећи одузимач_. Тренутно је то само `I` (1), али може бити и `X` (10) или `C` (100).
    - Ако задовољава оба ова услова, морамо да потражимо вредност и додамо је укупном _а_ ако је то један од посебних одузимача, у супротном занемарите
    - Тада треба даље повећавати `i` тако да овај симбол не рачунамо два пута


## Рефактор

Нисам у потпуности уверен да ће ово бити дугорочни приступ и потенцијално бисмо могли да направимо неке занимљиве рефакторе, али одупрћу се томе у случају да је наш приступ потпуно погрешан. Волео бих да прво направим још неколико тестова и видим. У међувремену сам прву изјаву `if` дао нешто мање ужасну.

```go
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		if couldBeSubtractive(i, symbol, roman) {
			nextSymbol := roman[i+1]

			// build the two character string
			potentialNumber := string([]byte{symbol, nextSymbol})

			// get the value of the two character string
			value := allRomanNumerals.ValueOf(potentialNumber)

			if value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++
			}
		} else {
			total++
		}
	}
	return total
}

func couldBeSubtractive(index int, currentSymbol uint8, roman string) bool {
	return index+1 < len(roman) && currentSymbol == 'I'
}
```

## Прво напишите тест

Пређимо на `cases[:5]`

```
=== RUN   TestConvertingToArabic/'V'_gets_converted_to_5
    --- FAIL: TestConvertingToArabic/'V'_gets_converted_to_5 (0.00s)
        numeral_test.go:62: got 1, want 5
```

## Напишите довољно кода да прође

Осим када је одузимајући, наш код претпоставља да је сваки знак `I`, због чега је вредност 1. Требали бисмо бити у могућности да поново користимо нашу методу `ValueOf` да бисмо то поправили.

```go
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		// look ahead to next symbol if we can and, the current symbol is base 10 (only valid subtractors)
		if couldBeSubtractive(i, symbol, roman) {
			nextSymbol := roman[i+1]

			// build the two character string
			potentialNumber := string([]byte{symbol, nextSymbol})

			if value := allRomanNumerals.ValueOf(potentialNumber); value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++ // this is fishy...
			}
		} else {
			total+=allRomanNumerals.ValueOf(string([]byte{symbol}))
		}
	}
	return total
}
```

## Рефактор

Када индексујете низове у програму Го, добијате `byte`. Због тога када поново градимо низ, морамо радити ствари попут `string([]byte{symbol})`. Понавља се неколико пута, хајде да само преместимо ту функционалност тако да уместо тога `ValueOf` узме неколико бајтова.

```go
func (r RomanNumerals) ValueOf(symbols ...byte) int {
	symbol := string(symbols)
	for _, s := range r {
		if s.Symbol == symbol {
			return s.Value
		}
	}

	return 0
}
```

Тада можемо само да пређемо у бајтовима као што је, у нашу функцију

```go
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		if couldBeSubtractive(i, symbol, roman) {
			if value := allRomanNumerals.ValueOf(symbol, roman[i+1]); value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++ // this is fishy...
			}
		} else {
			total+=allRomanNumerals.ValueOf(symbol)
		}
	}
	return total
}
```

Још увек је прилично гадно, али стиже тамо.

Ако почнете да премештате наш број `cases[:xx]`, видећете да их сада пролази доста. У потпуности уклоните оператор пресека и погледајте који не успевају, ево неколико примера из мог пакета

```
=== RUN   TestConvertingToArabic/'XL'_gets_converted_to_40
    --- FAIL: TestConvertingToArabic/'XL'_gets_converted_to_40 (0.00s)
        numeral_test.go:62: got 60, want 40
=== RUN   TestConvertingToArabic/'XLVII'_gets_converted_to_47
    --- FAIL: TestConvertingToArabic/'XLVII'_gets_converted_to_47 (0.00s)
        numeral_test.go:62: got 67, want 47
=== RUN   TestConvertingToArabic/'XLIX'_gets_converted_to_49
    --- FAIL: TestConvertingToArabic/'XLIX'_gets_converted_to_49 (0.00s)
        numeral_test.go:62: got 69, want 49
```

Мислим да нам недостаје само ажурирање `couldBeSubtractive` тако да узима у обзир остале врсте субтрактивних симбола

```go
func couldBeSubtractive(index int, currentSymbol uint8, roman string) bool {
	isSubtractiveSymbol := currentSymbol == 'I' || currentSymbol == 'X' || currentSymbol =='C'
	return index+1 < len(roman) && isSubtractiveSymbol
}
```

Покушајте поново, и даље не успевају. Међутим, оставили смо коментар раније ...

```go
total++ // this is fishy...
```

Никада не бисмо требали само повећавати `total`, јер то подразумева да је сваки симбол `I`. Замените га са:

```go
total += allRomanNumerals.ValueOf(symbol)
```

И сви тестови пролазе! Сад кад имамо потпуно радни софтвер, можемо се с поуздањем препустити некој рефакторизацији.

## Рефактор

Овде је сав код који сам завршио. Имао сам неколико неуспелих покушаја, али како стално наглашавам, то је у реду и тестови ми помажу да се слободно поигравам са кодом.

```go
import "strings"

func ConvertToArabic(roman string) (total int) {
	for _, symbols := range windowedRoman(roman).Symbols() {
		total += allRomanNumerals.ValueOf(symbols...)
	}
	return
}

func ConvertToRoman(arabic int) string {
	var result strings.Builder

	for _, numeral := range allRomanNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}

	return result.String()
}

type romanNumeral struct {
	Value  int
	Symbol string
}

type romanNumerals []romanNumeral

func (r romanNumerals) ValueOf(symbols ...byte) int {
	symbol := string(symbols)
	for _, s := range r {
		if s.Symbol == symbol {
			return s.Value
		}
	}

	return 0
}

func (r romanNumerals) Exists(symbols ...byte) bool {
	symbol := string(symbols)
	for _, s := range r {
		if s.Symbol == symbol {
			return true
		}
	}
	return false
}

var allRomanNumerals = romanNumerals{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

type windowedRoman string

func (w windowedRoman) Symbols() (symbols [][]byte) {
	for i := 0; i < len(w); i++ {
		symbol := w[i]
		notAtEnd := i+1 < len(w)

		if notAtEnd && isSubtractive(symbol) && allRomanNumerals.Exists(symbol, w[i+1]) {
			symbols = append(symbols, []byte{symbol, w[i+1]})
			i++
		} else {
			symbols = append(symbols, []byte{symbol})
		}
	}
	return
}

func isSubtractive(symbol uint8) bool {
	return symbol == 'I' || symbol == 'X' || symbol == 'C'
}
```

Мој главни проблем са претходним кодом је сличан нашем рефактору из ранијег. Имали смо превише брига повезаних заједно. Написали смо алгоритам који је покушавао да извуче римске бројеве из низа _и_, а затим пронашао њихове вредности.

Тако сам креирао нови тип `windowedRoman` који се побринуо за издвајање бројева, нудећи метод` Symbols` да их преузме као рез. То је значило да наша функција `ConvertToArabic` може једноставно прелазити преко симбола и збрајати их.

Код сам мало разградио извлачећи неке функције, посебно око изјаве винки иф да бих открио да ли је симбол којим тренутно имамо посла одузети симбол од два знака.

Вероватно постоји елегантнији начин, али нећу га знојити. Код је ту и делује и тестиран је. Ако ја (или било ко други) нађем бољи начин да могу безбедно да га промене - тежак посао је завршен.

## Увод у тестове засноване на својствима

У домену римских бројева било је неколико правила са којима смо радили у овом поглављу

- Не може имати више од 3 узастопна симбола
- Само I (1), X (10) и C (100) могу бити „одузимачи“
- Узимање резултата `ConvertToRoman(N)` и прослеђивање `ConvertToArabic` требало би да нам врати `N`

Тестови које смо до сада написали могу се описати као тестови засновани на „примерима“, где пружамо алатке за примере око нашег кода за верификацију.

Шта ако бисмо могли да прихватимо ова правила која знамо о свом домену и некако их применимо против свог кода?

Тестови засновани на својствима помажу вам у томе бацајући случајне податке на ваш код и верификујући да правила која описујете увек важе. Многи људи мисле да се тестови засновани на својствима углавном односе на случајне податке, али би погрешили. Прави изазов у тестовима заснованим на својствима је добро разумевање вашег домена како бисте могли да напишете ова својства.

Доста речи, да видимо неки код

```go
func TestPropertiesOfConversion(t *testing.T) {
	assertion := func(arabic int) bool {
		roman := ConvertToRoman(arabic)
		fromRoman := ConvertToArabic(roman)
		return fromRoman == arabic
	}

	if err := quick.Check(assertion, nil); err != nil {
		t.Error("failed checks", err)
	}
}
```

### Образложење вредности

Наш први тест ће проверити да ако трансформишемо број у римски, када користимо нашу другу функцију да бисмо га претворили у број да ћемо добити оно што смо првобитно имали.

- Дати случајни број (нпр. `4`).
- Позовите `ConvertToRoman` са случајним бројем (требало би вратити `IV` ако је `4`).
- Узмите горњи резултат и проследите га `ConvertToArabic`.
- Горе наведено требало би да нам пружи изворни унос (`4`).

Ово нам се чини као добар тест за изградњу самопоуздања, јер би требало да пукне ако у било којој грешци има грешке. Једини начин на који би то могло проћи је ако имају исту врсту грешке; што није немогуће, али се осећа мало вероватно.

### Техничко објашњење

Користимо пакет [testing/quick](https://golang.org/pkg/testing/quick/) из стандардне библиотеке

Читајући одоздо, пружамо функцију `quick.Check` која ће се покретати против одређеног броја случајних улаза, ако функција врати `false`, видеће се као неуспешна провера.

Наша горња функција `assertion` узима случајни број и покреће наше функције за тестирање својства.

### Покрените наш тест

Покушајте да га покренете; рачунар вам може неко време висити, па га убијте кад вам досади :)

Шта се дешава? Покушајте да додате следеће коду тврдње.


 ```go
assertion := func(arabic int) bool {
    if arabic <0 || arabic > 3999 {
        log.Println(arabic)
        return true
    }
    roman := ConvertToRoman(arabic)
    fromRoman := ConvertToArabic(roman)
    return fromRoman == arabic
}
```

Требали бисте видети нешто овако:

```
=== RUN   TestPropertiesOfConversion
2019/07/09 14:41:27 6849766357708982977
2019/07/09 14:41:27 -7028152357875163913
2019/07/09 14:41:27 -6752532134903680693
2019/07/09 14:41:27 4051793897228170080
2019/07/09 14:41:27 -1111868396280600429
2019/07/09 14:41:27 8851967058300421387
2019/07/09 14:41:27 562755830018219185
```

Само покретање овог врло једноставног својства открило је пропуст у нашој имплементацији. Као улаз користили смо `int`, али:
- Не можете да радите негативне бројеве са римским бројевима
- С обзиром на наше правило од највише 3 узастопна симбола, не можемо представити вредност већу од 3999 ([добро, некако](https://www.quora.com/Which-is-the-maximum-number-in-Roman-numerals)) и `int` има много већу максималну вредност од 3999.

Ово је супер! Били смо присиљени да дубље размислимо о свом домену који је стварна снага тестова заснованих на имовини.

Јасно је да `инт` није сјајан тип. Шта ако бисмо пробали нешто мало прикладније?


### [`uint16`](https://golang.org/pkg/builtin/#uint16)

Го има типове за _непотписане целе бројеве_, што значи да не могу бити негативни; тако да се одмах искључује једна класа грешака у нашем коду. Додавањем 16 то значи да је то 16-битни цели број који може да ускладишти максимум `65535`, што је и даље превелико, али нас приближава ономе што нам треба.

Покушајте да ажурирате код тако да користи `uint16` уместо `int`. Ажурирао сам `assertion` у тесту како бих пружио мало већу видљивост.

```go
assertion := func(arabic uint16) bool {
    if arabic > 3999 {
        return true
    }
    t.Log("testing", arabic)
    roman := ConvertToRoman(arabic)
    fromRoman := ConvertToArabic(roman)
    return fromRoman == arabic
}
```

Ако покренете тест, они сада стварно раде и можете видети шта се тестира. Можете покренути више пута да бисте видели да наш код добро стоји према различитим вредностима! Ово ми даје пуно самопоуздања да наш код ради како желимо.

Подразумевани број извођења `quick.Check` је 100, али то можете променити помоћу конфигурације.

```go
if err := quick.Check(assertion, &quick.Config{
    MaxCount:1000,
}); err != nil {
    t.Error("failed checks", err)
}
```

### Даљи рад

- Можете ли да напишете тестове својстава која проверавају друга својства која смо описали?
- Можете ли да смислите начин да то учините тако да је немогуће да неко зове наш код бројем већим од 3999?
    - Могли бисте да вратите грешку
    - Или креирајте нови тип који не може представљати> 3999
        - Шта мислиш да је најбоље?

## Окончање

### Још ТДД праксе са итеративним развојем

Да ли вам се помисао на писање кода који 1984. претвара у МЦМЛКСКСКСИВ у почетку осећала застрашујуће? Било ми је и већ дуго пишем софтвер.

Трик је, као и увек, **започети нешто једноставно** и предузети **мале кораке**.

Ни у једном тренутку у овом процесу нисмо направили неке велике скокове, направили било какве огромне преправке или ушли у неред.

Чујем како неко цинично говори „ово је само ката“. Не могу се расправљати с тим, али и даље користим исти приступ за сваки пројекат на којем радим. У првом кораку никада не испоручујем велики дистрибуирани систем, проналазим најједноставнију ствар коју би тим могао да испоручи (обично веб страницу „Хелло ворлд“), а затим понављам мале делове функционалности у управљачким деловима, баш као што смо то радили овде.

Вештина је знати _како_ поделити посао, а то долази са вежбом и неким дивним ТДД-ом који ће вам помоћи на путу.

### Испитивања заснована на својствима

- Уграђено у стандардну библиотеку
- Ако можете да смислите начине за описивање правила домена у коду, она су одличан алат који вам даје више самопоуздања
- Присилите вас да дубоко размислите о свом домену
- Потенцијално лепа допуна вашем тест пакету

## Postscript

Ова књига се ослања на драгоцене повратне информације из заједнице.
[Dave](http://github.com/gypsydave5)  је од огромне помоћи у практично сваком поглавље. Али он се стварно бунио о мојој употреби 'арапских бројева' у овоме поглавље, у интересу потпуног обелодањивања, ево шта је рекао.

> Само ћу написати зашто вредност типа `int` заправо није "арапски"
> број '. Ово је можда превише прецизно, па ћу у потпуности разумети
> ако ми кажеш да се искључим.
>
> _Дигит_ је знак који се користи за представљање бројева - од латинског
> за 'прст', као што их обично имамо десет. На арапском (такође се зове
> Хинду-арапски) бројевни систем има их десет. Ове арапске цифре су:
>
>     0 1 2 3 4 5 6 7 8 9
>
> _Број_ представља приказ броја помоћу збирке цифара.
> Арапски број је број представљен арапским цифрама у основи 10
> позициони систем бројева. Кажемо „позициони“ јер свака цифра има
> различита вредност на основу њеног положаја у бројци. Тако
>
>     1337
>
> `1` има вредност хиљаду јер је прва цифра у четворки
> цифрени број.
>
> Римљани се граде помоћу смањеног броја цифара (`I`, `V` итд ...) углавном као
> вредности за добијање броја. Има мало позиционих ствари, али то је
> углавном `/home/marcetin/GoProjects/nauci-go-sa-testovima/roman-numerals.md` увек представљам 'један'.
>
> Дакле, с обзиром на ово, да ли је `int` арапски број '? Идеја броја није у томе
> све везано за његову репрезентацију - то можемо видети ако се запитамо шта је
> тачан приказ овог броја је:
>
>     255
>     11111111
>     two-hundred and fifty-five
>     FF
>     377
>
> Да, ово је трик питање. Сви су тачни. Они су представништво
> истог броја у децималном, бинарном, енглеском, хексадецималном и окталном облику
> бројевни системи.
>
> Приказ броја као броја је _зависан_ од његових својстава
> као број - а то можемо видети када у Го-у погледамо целобројне литерале:
>
> ```go
>  0xFF == 255 // true
> ```
>
> И како можемо исписати читаве бројеве у низу формата:
>
> ```go
> n := 255
> fmt.Printf("%b %c %d %o %q %x %X %U", n, n, n, n, n, n, n, n)
> // 11111111 ÿ 255 377 'ÿ' ff FF U+00FF
> ```
>
> Можемо написати исти цијели број и као хексадецимални и као арапски (децимални)
> број.
>
> Дакле, када потпис функције изгледа као `ConvertToRoman(arabic int) string`
> помало претпоставља како се зове. Јер
> понекад ће се `arabic` писати као децимални цео број
>
> ```go
> ConvertToRoman(255)
> ```
>
> Али могло би се исто тако написати
>
> ```go
> ConvertToRoman(0xFF)
> ```
>
> Заиста, уопште не „конвертујемо“ из арапског броја, већ „штампамо“ -
> представљање - `int` као римски број - и `int` нису бројеви,
> Арапски или на неки други начин; то су само бројеви. Функција `ConvertToRoman` је
> више попут `strconv.Itoa` по томе што претвара `int` у `string`.
>
> Али сваку другу верзију кате није брига за ову разлику
> :shrug:
