# Читање датотека

- **[Сав код за ово поглавље можете пронаћи овде](https://github.com/marcetin/nauci-go-sa-testovima/tree/main/reading-files)**
- [Ево видео записа у којем радим кроз проблем узимајући питања из "Twitch stream"-а](https://www.youtube.com/watch?v=nXts4dEJnkU)

У овом поглављу ћемо научити како да прочитате неке датотеке, извуците неке податке и учините нешто корисно.

Претварајте се да радите са својим пријатељем да бисте креирали неки софтвер за блог. Идеја је аутор ће написати своје постове у ознаку, са неким метаподацима на врху датотеке. На стартупу, веб сервер ће прочитати мапу да створи неке `Post`-ове, а затим ће функција засебног `NewHandler`-а користиће оне `Post`-ове као датаСоурце за веб локацију блога.

Од мене су тражили да направимо пакет који претвара дату фасциклу Блог датотека у прикупљање `Post`-ова.

### Пример података

hello world.md
```markdown
Title: Hello, TDD world!
Description: First post on our wonderful blog
Tags: tdd, go
---
Hello world!

The body of posts starts after the `---`
```

### Очекивани подаци

```go
type Post struct {
	Title, Description, Body string
	Tags []string
}
```

## Итеративни, развојни развој

Узећемо итеративни приступ где увек узимамо једноставне, сигурне кораке према нашем циљу.

Ово захтева да прекинемо свој рад, али требали бисмо пазити да не паднемо у замку узимања ["од доле према горе"](https://en.wikipedia.org/wiki/Top-down_and_bottom-up_design) приступ.

Не треба да верујемо нашим активним маштема када започнемо са радом. Могли бисмо бити у искушењу да направимо неку врсту апстракције која је само потврђена када све будемо заједно, попут неке врсте `BlogPostFileParser`.

Ово је _not_ итеративно и недостаје се на уске подлоге повратних информација које ће нас ТДД довести.

Кент Бек каже:

> Оптимизам је опасност за професионал од програма. Повратне информације је третман.

Уместо тога, наш приступ треба да тежи да буде што је могуће ближе испоруци _реалних_ вредности конзумента што је брже могуће (често назива "срећним путем"). Једном када смо доставили малу количину потрошачке вредности крајње до краја, додатна итерација остатка захтева обично је једноставна.

## Размишљање о врсти теста који желимо да видимо

Подсетимо се нашег начина размишљања и циљева приликом почетка:

- **Напишите тест који желимо да видимо**. Размислите о томе како бисмо желели да користимо код који ћемо писати са становишта конзумента.
- Фокусирајте се на _шта_ и _зашто_, али немојте да вас омета _како_.

Наш пакет треба да понуди функцију која се може указати у мапу и врати нам неке постове.

```go
var posts []blogposts.Post
posts = blogposts.NewPostsFromFS("some-folder")
```

Да бисмо написали тест око тога, требаће нам нека врста тестне мапе са неким примерима у њему. _Нема ништа погрешно у реду с тим_, већ правите неке комбинације:

- За сваки тест можда ћете морати да креирате нове датотеке да бисте тестирали одређено понашање
- неко понашање ће бити изазовно за тестирање, као што је неуспех у учитавању датотека
- Тестови ће покренути мало спорије јер ће им требати приступити датотечном систему

Такође се непотребно не спојимо на специфичну примену датотечног система.

### Апстракције датотечног система уведене у Го 1.16

Го 1.16 је представио апстракцију за датотеке; [io/fs](https://golang.org/pkg/io/fs/) пакет.

> Пакет `fs` дефинише основне интерфејсе у датотечни систем. Систем датотека може да обезбеди оперативни систем домаћина, али и други пакети.

Ово нам омогућава да се отпустимо наше спојнице одређеном датотечном систему, који ће нам тада допустити различите имплементације у складу са нашим потребама.

> [На страни произвођача интерфејса, нови уградни.FS тип имплементира fs.FS, као и zip.Reader. Нова функција os.DirFS пружа имплементацију fs.FS подржаних од стабла датотека оперативног система.](https://golang.org/doc/go1.16#fs)

Ако користимо овај интерфејс, корисници нашег пакета имају низ опција печених у стандардној библиотеци која ће се користити. Учење улагања интерфејса дефинисаних у Го-овој стандардној библиотеци (нпр. `io.fs`, [`io.Reader`](https://golang.org/pkg/io/#Reader), [`io.Writer`](https://golang.org/pkg/io/#Writer)), је од виталног значаја за писање лаганих повезаних пакета. Ови пакети се затим могу поново користити у контекстима различитих од оних које сте замислили, уз минималну буку од ваших потрошача.

У нашем случају, можда наши потрошач жели да се положаје уграде у траје бинарну, а не датотеке у "стварном" датотечном систему? Било како било, _наш код не треба да брине о томе_.

За наше тестове, паковање [testing/fstest](https://golang.org/pkg/testing/fstest/) нуди нам имплементацију [io/FS](https://golang.org/pkg/io/fs/#FS) Да бисте користили, слично са алаткама које смо упознати у [net/http/httptest](https://golang.org/pkg/net/http/httptest/).

С обзиром на ове информације, следеће се осећа као бољи приступ,

```go
var posts blogposts.Post
posts = blogposts.NewPostsFromFS(someFS)
```


## Прво напишите тест

Треба да наставимо што је могуће мање и корисно што је могуће више. Ако докажемо да можемо прочитати све датотеке у директорију, то ће бити добар почетак. Ово ће нам дати поверење у софтвер који пишемо. Можемо да проверимо да је тачка враћања `[]Post` иста као и број датотека у нашем лажном систему датотека.

Креирајте нови пројекат за рад кроз ово поглавље.

- `mkdir blogposts`
- `cd blogposts`
- `go mod init github.com/{your-name}/blogposts`
- `touch blogposts_test.go`

```go
package blogposts_test

import (
	"testing"
	"testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {
    fs := fstest.MapFS{
        "hello world.md":  {Data: []byte("hi")},
        "hello-world2.md": {Data: []byte("hola")},
    }

    posts := blogposts.NewPostsFromFS(fs)

    if len(posts) != len(fs) {
        t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
    }
}

```

Примјетите да је пакет нашег теста `blogposts_test`. Запамтите, када се ТДД исправља добро, узмимо приступ _вођен конзументом_: не желимо да тестирамо интерне детаље јер _конзументе_ није брига за њих. Долажењем `_test` на наше намењено име пакета, приступимо само извозним члановима из нашег пакета - баш као и прави корисник нашег пакета.

Увезли смо [`testing/fstest`](https://golang.org/pkg/testing/fstest/) који нам даје приступ [`fstest.MapFS`](https://golang.org/pkg/testing/fstest/#MapFS) тип. Наш лажни систем датотека ће проћи `fstest.MapFS` у наш пакет.

> MapFS је једноставан систем датотека у меморији за употребу у тестовима, заступљен као мапа из имена стаза (аргументи за отварање) на информације о датотекама или именицима које представљају.

Ово је једноставније од одржавања мапе тестних датотека и то ће се брже извршити.

Коначно, кодификовали смо употребу нашег АПИ-ја са потрошачког становишта, а затим проверени да ли ствара тачан број постова.

## Покушајте да покренете тест

```
./blogpost_test.go:15:12: undefined: blogposts
```

## Напишите минималну количину кода за тест да бисте покренули и _проверили излазне податке неуспешног теста_

Паковање не постоји. Креирајте нову датотеку `blogposts.go` и ставите `package blogposts` у њему. Тада ћете морати да увезете тај пакет у своје тестове. За мене, увоз сада изгледа:


```go
import (
	blogposts "github.com/marcetin/nauci-go-sa-testovima/reading-files"
	"testing"
	"testing/fstest"
)
```

Сада се тестови неће компајлирати јер наш нови пакет нема функцију `NewPostsFromFS`, која враћа неку врсту колекције.

```
./blogpost_test.go:16:12: undefined: blogposts.NewPostsFromFS
```

Ово нас присиљава да направимо костур наших функција да тестирам. Не заборавите да у овом тренутку не претерано кванути кодекс; Покушавамо само да се покренемо тест и да се побринемо да не успе како бисмо очекивали. Ако прескочимо овај корак, можемо прескочити претпоставке и написати тест који није користан.

```go
package blogposts

import "testing/fstest"

type Post struct {

}

func NewPostsFromFS(fileSystem fstest.MapFS) []Post {
	return nil
}
```

Тест би сада требао исправно пасти

```
=== RUN   TestNewBlogPosts
    blogposts_test.go:48: got 0 posts, wanted 2 posts
```

## Напишите довољно кода да прође

We _could_ ["slime"](https://deniseyu.github.io/leveling-up-tdd/) this to make it pass:
Можемо урадити ["slime"](https://deniseyu.github.io/leveling-up-tdd/) на овом да би `прошло`:

```go
func NewPostsFromFS(fileSystem fstest.MapFS) []Post {
	return []Post{{},{}}
}
```

Али, како је Денисе Иу написао:

> Слиминг је корисно за давање "костура" вашем објекту. Дизајн интерфејса и извршења логике су две забринутости и тестирани тестови који се чине стратешки омогућава да се фокусирате на један по један по један.

Већ имамо своју структуру. Па, шта да радимо уместо тога?

Док смо смакли обим, све што требамо је да прочитамо директориј и креирамо пост за сваку датотеку коју сусрећемо. Не морамо да бринемо о отварању датотека и још увек их палимо.

```go
func NewPostsFromFS(fileSystem fstest.MapFS) []Post {
	dir, _ := fs.ReadDir(fileSystem, ".")
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts
}
```

[`fs.ReadDir`](https://golang.org/pkg/io/fs/#ReadDir) чита директоријум унутар дате `fs.FS` враћајући [`[]DirEntry`](https://golang.org/pkg/io/fs/#DirEntry).

Већ је наш идеализовани поглед на свет фолирао јер се могу догодити грешке, али се сећају да је сада наша фокус _направити тест да прође_, не мењајући дизајн, па ћемо за сада игнорисати грешку.

Остатак кода је једноставан: итерација преко уноса, направите `Post` за сваки и вратите се начек.

## Рефактор

Иако наши тестови пролазе, не можемо да користимо наш нови пакет изван овог контекста, јер је повезано са конкретном применом `fstest.MapFS`. Али то не мора бити. Промените аргумент на наше функције `NewPostsFromFS` да бисте прихватили интерфејс из стандардне библиотеке.

```go
func NewPostsFromFS(fileSystem fs.FS) []Post {
	dir, _ := fs.ReadDir(fileSystem, ".")
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts
}
```

Поновно покрените тестове: Све би требало да ради.

### Руковање грешком

Паркирали смо руковање грешком раније када смо се фокусирали на посао са радом срећног пута. Пре него што наставите са понављањем функционалности, требало би да признамо да се грешке могу догодити током рада са датотекама. Поред читања директорија, можемо наићи на проблеме када отворимо појединачне датотеке. Променимо наш АПИ (прво путем наших тестова прво, наравно) тако да може да врати `error`..

```go
func TestNewBlogPosts(t *testing.T) {
    fs := fstest.MapFS{
        "hello world.md":  {Data: []byte("hi")},
        "hello-world2.md": {Data: []byte("hola")},
    }

    posts, err := blogposts.NewPostsFromFS(fs)

    if err != nil {
        t.Fatal(err)
    }

    if len(posts) != len(fs) {
        t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
    }
}
```

Покрените тест: Треба се жалити на погрешан број повратних вредности. Поправка кода је једноставна.

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts, nil
}
```

Ово ће направити пролаз. ТДД практикант у вама можда је нервиран да нисмо видели неуспешни тест пре него што напишемо шифру да би пропагирали грешку из `fs.ReadDir`. Да би то требало "правилно", требали би нам нови тест где убризгамо неуспех `fs.FS` тест-двоструки да би направио `fs.ReadDir` Вратите `error`.

```go
type StubFailingFS struct {
}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}

// later
_, err := blogposts.NewPostsFromFS(StubFailingFS{})
```

Ово би требало да вам пружи самопоуздање у наш приступ. Интерфејс који користимо има једну методу, што ствара стварање тест-парова да би се тестирали различити сценарији тривијални.

У неким случајевима, тестирање руковања грешкама је прагматична ствар, али у нашем случају не радимо ништа _интересантно_ са грешком, само га пропагирамо, тако да не вреди гњаваже око писања новог теста.

Логично, наше следеће итерације ће се ширити на нашем типу `Post` да има неке корисне податке.

## првопишите тест

Почећемо са првом линијом у предложеном блогу поље, насловним пољем.

Морамо да променимо садржај тестних датотека тако да се подударају са оним што је наведен, а затим можемо да поднесемо тврдњу да је правилно рашчлањен.

```go
func TestNewBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte("Title: Post 1")},
		"hello-world2.md": {Data: []byte("Title: Post 2")},
	}

	// rest of test code cut for brevity
    got := posts[0]
    want := blogposts.Post{Title: "Post 1"}

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %+v, want %+v", got, want)
    }
}
```

## Покушајте да покренете тест
```
./blogpost_test.go:58:26: unknown field 'Title' in struct literal of type blogposts.Post
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

Додајте ново поље на наш `Post` тип тако да се тест покрене

```go
type Post struct {
	Title string
}
```

Поново покрените тест и требали бисте добити јасан, неуспешни тест

```
=== RUN   TestNewBlogPosts
=== RUN   TestNewBlogPosts/parses_the_post
    blogpost_test.go:61: got {Title:}, want {Title:Post 1}
```

## Напишите довољно кода да прође

Треба да отворимо сваку датотеку, а затим да издвојимо наслов

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(fileSystem, f)
		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
	postFile, err := fileSystem.Open(f.Name())
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()

	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

Запамтите да је наш фокус у овом тренутку не писати елегантног кода, то је само да би стигао до тачке у којој имамо радни софтвер.

Иако се то осећа као мали прираштај, ипак нам је ипак захтевао да напишемо фер кодекс и дајемо неке претпоставке у погледу руковања грешкама. Ово би била тачка у којој би требали разговарати са колегама и одлучити најбољи приступ.

Итеративни приступ нам је дао брзе повратне информације да је наше разумевање захтева непотпуно.

`fs.FS` нам даје начин да отворимо датотеку у њему по имену са његовим `Open` методом. Одатле смо прочитали податке из датотеке и за сада нам не треба ништа софистицирано рашчлањивање, само исећи `Title:` Текст резањем низа.

## Рефактор

Раздвајање `кода за отварање датотеке` са `кодом за читање садржаја датотека` учиниће код једноставнијим за разумевање и даљи рад.

```go
func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
	postFile, err := fileSystem.Open(f.Name())
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile fs.File) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

Када репродуцирате нове функције или методе, побрините се и размислите о аргументима. Овде дизајнирате и слободно је да дубоко размишљате о ономе што је прикладно јер имате пролазне тестове. Размислите о спојности и кохезији. У овом случају треба да се поставите:

> Да ли `newPost` мора бити спојен на `fs.File`? Да ли користимо све методе и податке са ове врсте? Шта нам _стварно_ треба?

У нашем случају га користимо само као аргумент `io.ReadAll`, који је потребан `io.Reader`. Дакле, требали бисмо отпустити спојницу у нашој функцији и затражити `io.Reader`.

```go
func newPost(postFile io.Reader) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

Можете да направите сличан аргумент за функцију нашег `getPost`, која има аргумент `fs.DirEntry`, али једноставно зове `Name()` да бисте добили име датотеке. Не треба нам све то; Хајде да се децоуле из тог типа и прође име датотеке кроз низ. Ево потпуно враћеног кода:

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(fileSystem, f.Name())
		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile io.Reader) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

Од сада, већина наших напора може бити уредно садржана у `newPost`. Забринутост отварања и понављања преко датотека врши се, а сада се можемо фокусирати на екстракцију података за `Post` тип. Иако технички неопходно, датотеке су леп начин да логично групне ствари повезане заједно, па сам премештала `Post` тип и `newPost` у нову датотеку `post.go`.

### Тест помагач

Требали бисмо се побринути и за наше тестове. Много ћемо извршити тврдње на `Posts`, тако да бисмо требали написати неки код да помогнемо у томе

```go
func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
```

```go
assertPost(t, posts[0], blogposts.Post{Title: "Post 1"})
```

## Прво напишите тест

Прошилимо наш тест даље да бисмо извукли следећи ред из датотеке, опис. Док се то не учини да се прође сада треба да се осећа угодно и познато.

```go
func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1`
		secondBody = `Title: Post 2
Description: Description 2`
	)

	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

    // rest of test code cut for brevity
    assertPost(t, posts[0], blogposts.Post{
        Title: "Post 1",
        Description: "Description 1",
    })

}
```

## Покушајте да покренете тест

```
./blogpost_test.go:47:58: unknown field 'Description' in struct literal of type blogposts.Post
```

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

Додајте ново поље у `Post`.

```go
type Post struct {
	Title       string
	Description string
}
```

The tests should now compile, and fail.
Тестови би сада требало да се компајлирају, и падну.

```
=== RUN   TestNewBlogPosts
    blogpost_test.go:47: got {Title:Post 1
        Description: Description 1 Description:}, want {Title:Post 1 Description:Description 1}
```

## Напишите довољно кода да прође

Стандардна библиотека има практичну библиотеку која ће вам помоћи да скенирате путем података, линијом према линији; [`bufio.Scanner`](https://golang.org/pkg/bufio/#Scanner)

> Скенер нуди погодан интерфејс за читање података као што је датотека нових линија са ограниченим линијама текста.

```go
func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	scanner.Scan()
	titleLine := scanner.Text()

	scanner.Scan()
	descriptionLine := scanner.Text()

	return Post{Title: titleLine[7:], Description: descriptionLine[13:]}, nil
}
```

Ручно је, такође је потребно `io.Reader` прочитати (хвала вам опет, лабав-спојница), не морамо да мењамо аргументе функције.

Позовите `Scan` да бисте прочитали линију, а затим извуците податке користећи `Text`.

Ова функција никада није могла да врати `error`. У овом тренутку би било примамљиво да га уклони из врсте повратка, али знамо да ћемо касније морати да се бавимо неважећим датотекама касније, тако да га можемо и оставити.

## Рефактор

Имамо понављање око скенирања линије и затим читајући текст. Ми знамо да ћемо радити ову операцију бар још једном, то је једноставан рефалтер да се осуши, па започнемо с тим.

```go
func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	title := readLine()[7:]
	description := readLine()[13:]

	return Post{Title: title, Description: description}, nil
}
```

То је једва спасило било које линије кода, али то је ретко тачка рефакторинга. Оно што овде покушавам да урадим само раздвајањем _шта_ од _како_ од редова за читање како би код учинио још мало декларативнијим за читаоца.

Иако чаробни бројеви од 7 и 13 обавите посао, нису ужасно описни.

```go
const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	title := readLine()[len(titleSeparator):]
	description := readLine()[len(descriptionSeparator):]

	return Post{Title: title, Description: description}, nil
}
```

Сада када зурим у код са својим креативним рефакторима, хтео бих да покушам да направим своју функцију Реадлине, побрините се за уклањање ознаке. Ту је и читљивији начин обријања префикса из низа са функцијом `strings.TrimPrefix`.

```go
func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Post{
		Title:       readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
	}, nil
}
```

Можете или не не свиђа ова идеја, али ја то радим. Поента је у стању рефакторинга, слободно се играмо са унутрашњим детаљима и можете наставити да користите своје тестове да проверите да се ствари и даље понашају правилно. Увек се можемо вратити на претходне државе ако нисмо задовољни. ТДД приступ нам ову лиценцу даје често експериментише са идејама, па имамо више снимака на писању сјајног кода.

Следећи захтев је екстрахирајући ознаке поста. Ако пратите заједно, препоручио бих да је покушате да је сами проведете пре него што прочитате. Сада бисте требали имати добар, итеративни ритам и осећати се самоуверено да бисте извукли следећи ред и анализирали податке.

За сажетост, нећу проћи кроз кораке ТДД-а, али ево теста са доданим ознакама.

```go
func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker`
	)

    // rest of test code cut for brevity
    assertPost(t, posts[0], blogposts.Post{
        Title:       "Post 1",
        Description: "Description 1",
        Tags:        []string{"tdd", "go"},
    })
}
```

Само се варате ако само копирате и залепите оно што ја пишем. Да бисмо били сигурни да смо сви на истој страници, ево мог кода који укључује екстракцију ознака.

```go
const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator        = "Tags: "
)

func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Post{
		Title:       readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
		Tags:        strings.Split(readMetaLine(tagsSeparator), ", "),
	}, nil
}
```

Надам се да овде нема изненађења. Успели смо да поново користимо `readMetaLine` да бисте добили следећу линију за ознаке, а затим их поделили користећи `strings.Split`.

Последња итерација на нашем срећном путу је да издвојите тело.

Ево подсетника на предложени формат датотеке.

```
Title: Hello, TDD world!
Description: First post on our wonderful blog
Tags: tdd, go
---
Hello world!

The body of posts starts after the `---`
```

Већ смо прочитали прве 3 редове. Затим морамо да прочитамо још једну линију, одбацимо га и тада остатак датотеке садржи тело.

## Прво напишите тест

Промените тестне податке да бисте имали сепаратор и тело са неколико нових линија да бисте проверили да ли ухватимо сав садржај.

```go
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
    )
```

Додајте нашу тврдњу као и остале

```go
	assertPost(t, posts[0], blogposts.Post{
        Title:       "Post 1",
        Description: "Description 1",
        Tags:        []string{"tdd", "go"},
        Body: `Hello
World`,
    })
```

## Покушајте да покренете тест

```
./blogpost_test.go:60:3: unknown field 'Body' in struct literal of type blogposts.Post
```

Као што смо очекивали.

## Напиши минималну количину кода за покретање теста и провери неуспешне резултате теста

Add `Body` to `Post` and the test should fail.
Додајте `Body` у `Post` и тест не би требало да падне.

```
=== RUN   TestNewBlogPosts
    blogposts_test.go:38: got {Title:Post 1 Description:Description 1 Tags:[tdd go] Body:}, want {Title:Post 1 Description:Description 1 Tags:[tdd go] Body:Hello
        World}
```

## Напишите довољно кода да прође

1. Скенирајте следећи ред да бисте игнорисали "---" сепаратор.
2. Наставите скенирање док не преостаје ништа за скенирање.

```go
func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	title := readMetaLine(titleSeparator)
	description := readMetaLine(descriptionSeparator)
	tags := strings.Split(readMetaLine(tagsSeparator), ", ")

	scanner.Scan() // ignore a line

	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}
	body := strings.TrimSuffix(buf.String(), "\n")

	return Post{
		Title:       title,
		Description: description,
		Tags:        tags,
		Body:        body,
	}, nil
}
```

- `scanner.Scan()` Враћа `bool`, што указује да ли постоји још података за скенирање, тако да можемо да користимо то са `for` петљу да наставим читати кроз податке до краја.
- Након сваког `Scan()` Подаци пишемо у међуспремник користећи `fmt.Fprintln`. Ми користимо верзију која додаје нову линију јер скенер уклања нове линије из сваког ретка, али морамо их одржати.
- Због горе наведеног, морамо да подрегнемо коначну нову линију, тако да немамо задња.

## Рефактор

Капсулишући идеју да стављајући остатак података у функцију помаже будућим читаоцима да брзо разумеју `шта` се дешава у `newPost`, без да се брину о специфичностима за имплементацију.

```go
func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Post{
		Title:       readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
		Tags:        strings.Split(readMetaLine(tagsSeparator), ", "),
		Body:        readBody(scanner),
	}, nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan() // ignore a line
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}
	return strings.TrimSuffix(buf.String(), "\n")
}
```

## Иертирајући даље

Направили смо нашу "челичну нит" функционалности, узимајући најкраћу руту да стигнемо до наше срећне стазе, али очигледно да постоји нека удаљеност да крене пре него што је производња спремна.

Нисмо се обрадили:

- Када формат датотеке није тачан
- Датотека није `.md`
- Шта ако је редослед метаподатака другачије? Да ли би то требало дозволити? Да ли бисмо могли да то решимо?

Иако је пресудно, имамо радни софтвер и дефинисали смо наш интерфејс. Горе су само даље итерације, више тестова за писање и вођење нашег понашања. Да би подржали било који од горе наведеног, не бисмо требали да мењамо наш _дизајн_, већ само детаље о имплементацији.

Вођењи фокусом на циљ значи да смо донели важне одлуке и потврдили их против жељеног понашања, а не да се оптерећујемо питањима која неће утицати на целокупни дизајн.

## Окончање

`fs.FS`, а остале промене у Го 1.16 Дајте нам неке елегантне начине читања података из датотечних система и једноставно их тестирамо.

Ако желите да испробате код "за стварни":

- Креирајте мапу `cmd` у оквиру пројекта, додајте датотеку `main.go`
- Додајте следећи код

```go
import (
    blogposts "github.com/quii/fstest-spike"
    "log"
    "os"
)

func main() {
	posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(posts)
}
```

- Додајте неке датотеке ознаке у фасциклу `posts` и покрените програм!

Приметите симетрију између производног кода

```go
posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
```

И тестова

```go
posts, err := blogposts.NewPostsFromFS(fs)
```


Ово је када је управљање потрошачем, одоздо ТДД _делује исправно_.

Корисник нашег пакета може да погледа наше тестове и брзо устаје у брзини са оним што би требало да уради и како да га користи. Као одржани, можемо бити _уверени да су наши тестови корисни јер су из тренутне тачке конзумента_. Не тестирамо детаље о имплементацији или другим случајним детаљима, тако да можемо бити разумно сигурни да ће нам наши тестови помоћи, а не да нас ометају приликом рефакторинга.

Ослањајући се на добром софтверским инжењерским праксама попут [**убацивање пакета од којих зависи апликација**](dependency-injection.md) Наш код је једноставан за тестирање и поновну употребу.

Када креирате пакете, чак и ако су само унутрашњи за ваш пројекат, преферирајте приступ потрошача одоздо према доле. Ово ће вас зауставити претерано замишљање дизајна и прављење апстракција Можда чак и нећете потребан и помоћи ће вам да осигурате да су тестови које пишу корисни.

Итеративни приступ је чувао сваки ситни корак, а континуиране повратне информације помогле су нам да откријемо нејасне захтеве вероватно пре него што је остало, више ад-хоц приступи.

### Writing?

Важно је напоменути да ове нове функције имају само операције за датотеке _читање_. Ако ваш рад треба да пише, мораћете да погледате негде другде. Не заборавите да размишљате о томе шта тренутно нуди стандардна библиотека, ако пишете податке, вероватно бисте требали да потражите да искористите постојеће сучеље као што су `io.Writer` да бисте код кода задржели лабавци и поново употребљивали.

### Додатна литература

- Ово је било светло увод у `io/fs`. [Бен Цонгдон је урадио одличан текст](https://benjamincongdon.me/blog/2021/01/21/A-Tour-of-Go-116s-iofs-package/) Који је био пуно помоћи Писање овог поглавља.
- [Расправа о интерфејсима датотека](https://github.com/golang/go/issues/41190)
