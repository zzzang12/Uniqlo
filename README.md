# Uniqlo
유니클로 UT 크롤링 프로젝트

## 배운 것
* http GET 요청 보내기

  ```go
  func httpGet(url string) *http.Response {
    resp, err := http.Get(url)
    checkError(err)
    checkStatusCode(resp)
    return resp
  }
  ```
  이런 식으로 response를 받아온다.
  ```go
  resp := httpGet(url)
  defer resp.Body.Close()
  ```
  함수가 끝나면 ```Close```까지 해주기
  ```go
  func checkStatusCode(resp *http.Response) {
    if resp.StatusCode != 200 {
      log.Fatalf("status code error: %s", resp.Status)
    }
  }
  ```
  response를 잘 받았는 지도 확인해주기


* goquery로 html에서 원하는 node 가져오기

  ```go
  doc, err := goquery.NewDocumentFromReader(resp.Body)
  checkError(err)
  ```
  이렇게 response로부터 html document를 가져오고
  ```go
  sel := doc.Find("~~~")
  
  sel.Find("~~~").Each(func(i int, elem *goquery.Selection) {
    elem.Find("~~~").Attr("~~~")
    elem.Find("~~~").Text()
  })
  ```
  ```Find``` method는 receiver와 return value의 type이 모두 ```*goquery.Selection```이므로 찾아놓은 결과 값에 또 ```Find``` method를 이용하여 원하는 결과를 찾을 수 있다.

  또한 ```goquery.Selection```은 하나의 값을 가질 수도 있고 여러 값을 가질 수도 있어서 ```Find``` method를 통해 찾은 결과가 하나라면 하나의 값만 반환되지만 결과가 여러 개라면 여러 값이 반환된다.

  반환 값이 여러 개라면 ```Each```method에 익명함수를 넣어서 여러 반환 값에 대해 어떠한 처리를 해줄 수 있다.


* css selector
  
  ```Find``` method는 argument로 string type의 css selector가 들어간다.
  ```html
    <div class="thisIsClass" id="thisIsId">
      <h1>111</h1>
      <h1>222</h1>
      <h1>333</h1>
      <h2>444</h2>
      <h3>555</h3>
      <h3>666</h3>
      <h3>777</h3>
      <h4>888</h4>
      <h4>999</h4>
    </div>
    ```
    위 html을 예시로 들면
  * html 요소 선택자
    
    ```go
    sel := doc.Find("div")
    ```
    html 요소 선택자를 넣어서 html 태그를 검색
    
    ```sel```에 ```div```노드가 담김
  * class 선택
    ```go
    sel := doc.Find(".thisIsClass")
    ```
    앞에 ```.```을 붙여서 class 선택자를 넣어서 class를 검색
    
    ```sel```에 class가 ```thisIsClass```인 ```div```노드가 담김
  * id 선택자
    ```go
    sel := doc.Find("#thisIsId")
    ```
    앞에 ```#```을 붙여서 id 선택자를 넣어서 id를 검색

    ```sel```에 id가 ```thisIsId```인 ```div```노드가 담김
  * 자식 선택자
    ```go
    sel := doc.Find("div>h1")
    ```
    ```>```를 중간에 넣어서 가장 먼저 있는 자식 요소 하나를 검색

    ```sel```에 ```div```의 자식들 중 가장 위에 있는 ```h1```노드(111)가 담김
  * 자손 선택자
    ```go
    sel := doc.Find("div h1")
    ```
    ``` ```를 중간에 넣어서 모든 자식 요소를 검색

    ```sel```에 ```div```의 자식들 중 ```h1```태그를 가진 노드 모두(111, 222, 333)가 담김
  * 인접 형제 선택자
    ```go
    sel := doc.Find("h2+h3")
    ```
    ```+```를 중간에 넣어서 바로 다음 형제 요소를 검색

    ```sel```에 ```h2```의 바로 다음 형제인 ```h3```노드(555)가 담김
  * 일반 형제 선택자
    ```go
    sel := doc.Find("h2~h3")
    ```
    ```~```를 중간에 넣어서 모든 다음 형제 요소를 검색

    ```sel```에 ```h2```의 다음 형제들 중 ```h3```태그를 가진 노드 모두(555, 666, 777)가 담김
  

* directory 만들기

  ```go
  name := "folder"
  err := os.Mkdir(name, 0777)
  checkError(err)
  ```
  ```Mkdir``` method는 첫번째 argument로 만들어질 directory 이름, 두번째 argument로 권한을 주면 현재 위치에서 그 이름을 가진 directory를 만든다.
  ```go
  path := "folder" + "/" + "here"
  err := os.MkdirAll(path, 0777)
  checkError(err)
  ```
  ```MkdirAll``` method는 첫번째 argument로 만들어질 directory 경로를 주는데 ```Mkdir```과 달리 중첩된 directory를 만들 수 있다.


* file 만들기

  ```go
  path := "folder" + "/" + "file" + ".jpg"
  file, err := os.Create(path)
  checkError(err)
  defer file.Close()
  ```
  만들어질 file의 경로를 argument로 줘서 ```Create``` method를 통해 file을 만들 수 있고 file handler를 반환한다.

  함수가 끝나면 ```Close```까지 해주기


* image download
  
  image download는 file 하나를 만들고 거기에 image를 붙여넣는 방식으로 한다.

  ```go
  written, err := io.Copy(file, resp.Body)
  checkError(err)
  ```
  ```Copy``` method에 file handler와 http GET 요청으로 받아온 해당 image 주소의 response을 argument로 주면 file에 image가 담긴다.


* string parsing

  string을 parsing하는 방법은 여러가지가 있다.
  1. ```Split```
      ```go
      result := strings.Split("a,b,c", ",")
      result == ["a" "b" "c"]
      ```
      그냥 separator를 기준으로 분리해서 분리된 []string을 반환한다.
  2. ```SplitAfter```
      ```go
      result := strings.SplitAfter("a,b,c", ",")
      result == ["a," "b," "c"]
      ```
      spearator를 기준으로 분리하긴 하지만 seperator가 사라지진 않고 그까지 포함하고 다음 문자부터 분리한다.
  3. ```Fields```
      ```go
      result := strings.Fields("This is a string containing whitespaces")
      result == ["This" "is" "a" "string" "containing" "whitespace"]
      ```
      공백을 기준으로 간단하게 분리한다.
  4. ```FieldsFunc```
      ```go
      func split(c rune) bool {
        return !unicode.IsLetter(c) && !unicode.IsNumber(c)
      }
      result := strings.FieldsFunc("  foo1;bar2,baz3...", split)
      result == ["foo1" "bar2" "baz3"]
      ```
      다른 method들은 separator를 하나만 가지지만 ```FiledsFunc```는 여러 개를 가질 수 있다.
      
      두번째 arguent로 함수를 받는데, 이 함수는 separator로 쓰고 싶은 조건을 bool type으로 반환하면 그에 맞게 분리해준다.
  5. ```regexp```
  
      regex는 나중에 적절하게 사용할 때가 되면 그때 알아봐야겠다.
  

## 사용법
```shell
git clone github.com/zzzang12/uniqlo.git
go build uniqlo.go
./uniqlo.exe
```

## 시연 영상
<img src="https://user-images.githubusercontent.com/70265177/189478007-1583af6e-874f-427c-9aa8-0ff5d8d4cd30.gif" alt="video">