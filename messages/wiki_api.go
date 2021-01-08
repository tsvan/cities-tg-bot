package messages

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

//Конвертируем запрос для использование в качестве части URL
func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

//Result Структура ответа апи
type Result struct {
	Name, Description, URL string
}

//SearchResults тело респонса апи вики
type SearchResults struct {
	ready   bool
	Query   string
	Results []Result
}

//UnmarshalJSON - десериализация json-a
func (sr *SearchResults) UnmarshalJSON(bs []byte) error {
	array := []interface{}{}
	if err := json.Unmarshal(bs, &array); err != nil {
		return err
	}
	sr.Query = array[0].(string)
	for i := range array[1].([]interface{}) {
		sr.Results = append(sr.Results, Result{
			array[1].([]interface{})[i].(string),
			array[2].([]interface{})[i].(string),
			array[3].([]interface{})[i].(string),
		})
	}
	return nil
}

//GetWikiLink - получение ссылки по названию города
func GetWikiLink(request string) string {
	param, _ := urlEncoded(request)

	response, err := http.Get("https://ru.wikipedia.org/w/api.php?action=opensearch&search=" + param + "&limit=1&origin=*&format=json")
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	sr := &SearchResults{}
	if err = json.Unmarshal([]byte(contents), sr); err != nil {
		return ""
	}
	for i := range sr.Results {
		return sr.Results[i].URL
	}

	return ""
}
