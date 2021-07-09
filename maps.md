# Мапе

**[Сав код за ово поглавље можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/maps)**

У [низови и резови](arrays-and-slices.md), видели сте како сместити вредности у ред. Сада ћемо размотрити начин складиштења предмета помоћу `key` и брзо их потражити.

Мапе вам омогућавају да предмете складиштите на начин сличан речнику. `key` можете сматрати речју, а вредност „дефиницијом“. А који бољи начин постоји за учење о Мапама од стварања сопственог речника?

Прво, под претпоставком да већ имамо неке речи са њиховим дефиницијама у речнику, ако тражимо реч, она би требало да је врати у дефиницију.

## Прво напишите тест

У `dictionary_test.go`

```go
package main

import "testing"

func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    if got != want {
        t.Errorf("got %q want %q given, %q", got, want, "test")
    }
}
```

Декларирање мапе је донекле слично низу. Осим што започиње са кључном речи `map` и захтева два типа. Први је тип кључа који је написан унутар `[]`. Други је тип вредности, који иде одмах иза `[]`.

Тип кључа је посебан. То може бити упоредив тип, јер без могућности да се утврди да ли су 2 кључа једнака, не можемо да осигурамо да добијемо тачну вредност. Упоредиви типови детаљно су објашњени у [спецификацији језика](https://golang.org/ref/spec#Comparison_operators).

Тип вредности, с друге стране, може бити било који тип који желите. То може бити и друга карта.

Све остало у овом тесту би требало да буде познато.


## Покушајте да покренете тест

Покретањем `go test` компајлер неће успети са `./dictionary_test.go:8:9: undefined: Search`.

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

У `dictionary.go`

```go
package main

func Search(dictionary map[string]string, word string) string {
    return ""
}
```

Ваш тест би сада требало да пропадне са *јасном поруком о грешци*

`dictionary_test.go:12: got '' want 'this is just a test' given, 'test'`.

## Напишите довољно кода да прође

```go
func Search(dictionary map[string]string, word string) string {
    return dictionary[word]
}
```

Извлачење вредности из мапе исто је што и извлачење вредности из низа `map[key]`.

## Рефактор

```go
func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    assertStrings(t, got, want)
}

func assertStrings(t testing.TB, got, want string) {
    t.Helper()

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Одлучио сам да створим помоћника `assertStrings` како бих имплементацију учинио општијом.

### Коришћење прилагођеног типа

Употребу нашег речника можемо побољшати стварањем новог типа око мапе и начином `Search`.

У `dictionary_test.go`:

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    got := dictionary.Search("test")
    want := "this is just a test"

    assertStrings(t, got, want)
}
```

Почели смо да користимо тип `Dictionary`, који још увек нисмо дефинисали. Позовимо `Search` на `Dictionary` инстанци.

Нисмо морали да мењамо `assertStrings`.

У `dictionary.go`:

```go
type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
    return d[word]
}
```

Овде смо створили тип `Dictionary` који делује као танак омотач око мапе `map`. Са дефинисаним прилагођеним типом, можемо створити методу `Search`.

## Прво напишите тест

Основну претрагу било је врло лако спровести, али шта ће се догодити ако унесемо реч која није у нашем речнику?

Ми заправо ништа не враћамо. То је добро јер програм може да се изводи и даље, али постоји бољи приступ. Функција може извести да речи нема у речнику. На овај начин, корисник се не пита да ли та реч не постоји или једноставно нема дефиниције (ово се можда не чини врло корисним за речник. Међутим, то је сценарио који би могао бити кључан у другим случајевима коришћења).

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    t.Run("known word", func(t *testing.T) {
        got, _ := dictionary.Search("test")
        want := "this is just a test"

        assertStrings(t, got, want)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dictionary.Search("unknown")
        want := "could not find the word you were looking for"

        if err == nil {
            t.Fatal("expected to get an error.")
        }

        assertStrings(t, err.Error(), want)
    })
}
```

Начин за руковање овим сценаријем у Го-у је враћање другог аргумента типа `Error`.

Грешке се могу претворити у низ методом `.Error()`, што радимо када их проследимо тврдњи. Такође штитимо `assertStrings` са` if` да бисмо осигурали да не позивамо `.Error()` на `nil`.

## Покушајте да покренете тест

Ово не може да се компајлира

```
./dictionary_test.go:18:10: assignment mismatch: 2 variables but 1 values
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

```go
func (d Dictionary) Search(word string) (string, error) {
    return d[word], nil
}
```

Ваш тест би сада требало да пропадне са много јаснијом поруком о грешци.

`dictionary_test.go:22: expected to get an error.`

## Напишите довољно кода да прође

```go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", errors.New("could not find the word you were looking for")
    }

    return definition, nil
}
```

Да бисмо ово прошли, користимо занимљиво својство претраживања мапе. Може да врати 2 вредности. Друга вредност је логичка вредност која показује да ли је кључ пронађен успешно.

Ово својство нам омогућава да разликујемо реч која не постоји и реч која једноставно нема дефиницију.

## Рефактор

```go
var ErrNotFound = errors.New("could not find the word you were looking for")

func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", ErrNotFound
    }

    return definition, nil
}
```

Чаробне грешке у нашој функцији `Search` можемо се ослободити издвајањем у променљиву. Ово ће нам такође омогућити бољи тест.

```go
t.Run("unknown word", func(t *testing.T) {
    _, got := dictionary.Search("unknown")

    assertError(t, got, ErrNotFound)
})

func assertError(t testing.TB, got, want error) {
    t.Helper()

    if got != want {
        t.Errorf("got error %q want %q", got, want)
    }
}
```

Стварањем новог помагача успели смо да поједноставимо тест и почнемо да користимо променљиву `ErrNotFound` како наш тест не би пропао ако у будућности променимо текст грешке.

## Прво напишите тест

Имамо одличан начин претраживања речника. Међутим, ми не можемо да додамо нове речи у наш речник.

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    dictionary.Add("test", "this is just a test")

    want := "this is just a test"
    got, err := dictionary.Search("test")
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

У овом тесту користимо нашу функцију `Search` да бисмо мало олакшали валидацију речника.

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

У `dictionary.go`

```go
func (d Dictionary) Add(word, definition string) {
}
```

Ваш тест би сада требало да пропадне

```
dictionary_test.go:31: should find added word: could not find the word you were looking for
```

## Напишите довољно кода да прође

```go
func (d Dictionary) Add(word, definition string) {
    d[word] = definition
}
```

Додавање на мапу је такође слично низу. Потребно је само да наведете кључ и поставите га једнаким вредности.

### Показивачи, копије и др

Занимљиво својство мапа је да их можете модификовати без прослеђивања као адресе на њих (нпр. `&myMap`)

Ово их може учинити да се понашају као "референтни тип", [али како описује Даве Цхенеи](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it) нису.

> Вредност мапе је показивач на структуру runtime.hmap.

Дакле, када проследите мапу функцији / методи, заиста је копирате, али само део показивача, а не основну структуру података која садржи податке.

Каза са мапама је да оне могу имати вредност `nil`. Мапа `nil` понаша се као празна мапа током читања, али покушаји писања на мапу `nil` изазваће панику у току извршавања. Више о мапама можете прочитати [овде](https://blog.golang.org/go-maps-in-action).

Због тога никада не треба иницијализовати празну променљиву мапе:

```go
var m map[string]string
```

Уместо тога, можете иницијализовати празну мапу као што смо радили горе, или помоћу кључне речи `make` створити мапу за вас:

```go
var dictionary = map[string]string{}

// OR

var dictionary = make(map[string]string)
```

Оба приступа стварају празну `hash map` и усмеравају `dictionary` на њу. Што осигурава да никада нећете добити панику током рада.

## Рефактор

У нашој имплементацији нема много тога за рефакторирање, али тест би могао да искористи мало поједностављења.

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    word := "test"
    definition := "this is just a test"

    dictionary.Add(word, definition)

    assertDefinition(t, dictionary, word, definition)
}

func assertDefinition(t testing.TB, dictionary Dictionary, word, definition string) {
    t.Helper()

    got, err := dictionary.Search(word)
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if definition != got {
        t.Errorf("got %q want %q", got, definition)
    }
}
```

Направили смо променљиве за реч и дефиницију и преместили тврдњу дефиниције у сопствену помоћну функцију.

Наш `Add` изгледа добро. Осим тога, нисмо разматрали шта се дешава када вредност коју покушавамо да додамо већ постоји!

Мапа неће појавити грешку ако вредност већ постоји. Уместо тога, они ће наставити и преписати вредност са ново пруженом вредношћу. Ово може бити згодно у пракси, али име наше функције чини мање него тачним. `Add` не би требало да мења постојеће вредности. Требало би само да дода нове речи у наш речник.

## Прво напишите тест

```go
func TestAdd(t *testing.T) {
    t.Run("new word", func(t *testing.T) {
        dictionary := Dictionary{}
        word := "test"
        definition := "this is just a test"

        err := dictionary.Add(word, definition)

        assertError(t, err, nil)
        assertDefinition(t, dictionary, word, definition)
    })

    t.Run("existing word", func(t *testing.T) {
        word := "test"
        definition := "this is just a test"
        dictionary := Dictionary{word: definition}
        err := dictionary.Add(word, "new test")

        assertError(t, err, ErrWordExists)
        assertDefinition(t, dictionary, word, definition)
    })
}
...
func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

За овај тест смо модификовали `Add` да бисмо вратили грешку коју проверавамо у односу на нову променљиву грешке, `ErrWordExists`. Такође смо модификовали претходни тест да бисмо проверили да ли постоји грешка `nil`, као и функција` assertError`.

## Покушајте да покренете тест

Преводник неће успети јер не враћамо вредност за `Add`.

```
./dictionary_test.go:30:13: dictionary.Add(word, definition) used as value
./dictionary_test.go:41:13: dictionary.Add(word, "new test") used as value
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

У `dictionary.go`

```go
var (
    ErrNotFound   = errors.New("could not find the word you were looking for")
    ErrWordExists = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
    d[word] = definition
    return nil
}
```

Сада добијамо још две грешке. Још увек модификујемо вредност и враћамо грешку `nil`.

```
dictionary_test.go:43: got error '%!q(<nil>)' want 'cannot add word because it already exists'
dictionary_test.go:44: got 'new test' want 'this is just a test'
```

## Напишите довољно кода да прође

```go
func (d Dictionary) Add(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        d[word] = definition
    case nil:
        return ErrWordExists
    default:
        return err
    }

    return nil
}
```

Овде користимо израз `switch` да бисмо се подударали са грешком. Имати овакав `switch` пружа додатну сигурносну мрежу, у случају да` Search` врати грешку која није `ErrNotFound`.

## Рефактор

Немамо превише за рефакторирање, али како наша употреба грешака расте можемо направити неколико модификација.

```go
const (
    ErrNotFound   = DictionaryErr("could not find the word you were looking for")
    ErrWordExists = DictionaryErr("cannot add word because it already exists")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}
```

We made the errors constant; this required us to create our own `DictionaryErr` type which implements the `error` interface. You can read more about the details in [this excellent article by Dave Cheney](https://dave.cheney.net/2016/04/07/constant-errors). Simply put, it makes the errors more reusable and immutable.

Next, let's create a function to `Update` the definition of a word.

## Прво напишите тест

```go
func TestUpdate(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    dictionary.Update(word, newDefinition)

    assertDefinition(t, dictionary, word, newDefinition)
}
```

`Update` is very closely related to `Add` and will be our next implementation.

## Покушајте да покренете тест

```
./dictionary_test.go:53:2: dictionary.Update undefined (type Dictionary has no field or method Update)
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

We already know how to deal with an error like this. We need to define our function.

```go
func (d Dictionary) Update(word, definition string) {}
```

With that in place, we are able to see that we need to change the definition of the word.

```
dictionary_test.go:55: got 'this is just a test' want 'new definition'
```

## Напишите довољно кода да прође

We already saw how to do this when we fixed the issue with `Add`. So let's implement something really similar to `Add`.

```go
func (d Dictionary) Update(word, definition string) {
    d[word] = definition
}
```

There is no refactoring we need to do on this since it was a simple change. However, we now have the same issue as with `Add`. If we pass in a new word, `Update` will add it to the dictionary.

## Прво напишите тест

```go
t.Run("existing word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    newDefinition := "new definition"
    dictionary := Dictionary{word: definition}

    err := dictionary.Update(word, newDefinition)

    assertError(t, err, nil)
    assertDefinition(t, dictionary, word, newDefinition)
})

t.Run("new word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{}

    err := dictionary.Update(word, definition)

    assertError(t, err, ErrWordDoesNotExist)
})
```

We added yet another error type for when the word does not exist. We also modified `Update` to return an `error` value.

## Покушајте да покренете тест

```
./dictionary_test.go:53:16: dictionary.Update(word, "new test") used as value
./dictionary_test.go:64:16: dictionary.Update(word, definition) used as value
./dictionary_test.go:66:23: undefined: ErrWordDoesNotExist
```

We get 3 errors this time, but we know how to deal with these.

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

```go
const (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

func (d Dictionary) Update(word, definition string) error {
    d[word] = definition
    return nil
}
```

We added our own error type and are returning a `nil` error.

With these changes, we now get a very clear error:

```
dictionary_test.go:66: got error '%!q(<nil>)' want 'cannot update word because it does not exist'
```

## Напишите довољно кода да прође

```go
func (d Dictionary) Update(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        d[word] = definition
    default:
        return err
    }

    return nil
}
```

Ова функција изгледа готово идентично `Add`, осим што смо се пребацили када ажурирамо `dictionary` и када вратимо грешку.

### Напомена о пријављивању нове грешке за Ажурирање

Могли бисмо поново да користимо `ErrNotFound` и да не додамо нову грешку. Међутим, често је боље имати прецизну грешку када ажурирање не успе.

Ако имате одређене грешке, добићете више информација о томе шта је пошло по злу. Ево примера у веб апликацији:

> Можете преусмерити корисника када се наиђе на `ErrNotFound`, али приказати поруку о грешци када се наиђе на` ErrWordDoesNotExist`.

Даље, креирајмо функцију за `Delete` речи из речника.

## Прво напишите тест

```go
func TestDelete(t *testing.T) {
    word := "test"
    dictionary := Dictionary{word: "test definition"}

    dictionary.Delete(word)

    _, err := dictionary.Search(word)
    if err != ErrNotFound {
        t.Errorf("Expected %q to be deleted", word)
    }
}
```

Наш тест креира `Dictionary` са речју, а затим проверава да ли је реч уклоњена.

## Покушајте да покренете тест

Покретањем `go test` добијамо:

```
./dictionary_test.go:74:6: dictionary.Delete undefined (type Dictionary has no field or method Delete)
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

```go
func (d Dictionary) Delete(word string) {

}
```

Након што ово додамо, тест нам говори да реч не бришемо.

```
dictionary_test.go:78: Expected 'test' to be deleted
```

## Напишите довољно кода да прође

```go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

Го има уграђену функцију `delete` која ради на мапама. Потребна су два аргумента. Прва је мапа, а друга је кључ који треба уклонити.

Функција `delete` не враћа ништа, а ми смо нашу методу` Delete` засновали на истом појму. Будући да брисање вредности која не постоји нема ефекта, за разлику од наших метода `Update` и `Add`, не треба да компликујемо АПИ грешкама.

## Окончање

У овом одељку смо доста обрадили. Направили смо пуни CRUD (Create, Read, Update and Delete) АПИ за наш речник. Током процеса научили смо како:

* Креирајте мапе
* Потражите ставке на мапама
* Додајте нове ставке на мапе
* Ажурирајте ставке на мапама
* Избришите ставке са мапе
* Сазнајте више о грешкама
    * Како створити грешке које су константе
    * Писање омота са грешкама
