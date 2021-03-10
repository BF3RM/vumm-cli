package registry

import (
	"github.com/vumm/cli/common"
)

func GetMod(mod string) (res Mod, err error) {
	err = common.GetHttpJson(Url+"/mods/"+mod, &res)
	return
}
