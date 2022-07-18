package elasticconvert

import (
	"encoding/json"
	"fmt"

	"github.com/suiguo/elastic-convert/model"
)

var ErrorParamsLen error = fmt.Errorf("")

//data es resp    len(indexs) = len(targets)
func Result(data []byte, indexs []string, targets []interface{}) (out *model.Result, err error) {
	out = &model.Result{}
	if len(indexs) != len(targets) {
		return out, ErrorParamsLen
	}
	for idx, index := range indexs {
		model.ModelMap.Store(index, targets[idx])
	}
	err = json.Unmarshal(data, out)
	return out, err
}
