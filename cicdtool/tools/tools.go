package tools

import (
    "fmt"
    "strconv"
    "strings"
)

//v0.0.0
func GetNewVersion(version string)string{
    versionNV := strings.Split(version, "v")
    versionNum := strings.Split(versionNV[1], ".")
    for i := len(versionNum)-1; i >=0 ;i-- {
        num, _ := strconv.Atoi(versionNum[i])
        if num<9{
            num++
            versionNum[i]=fmt.Sprintf("%d",num)
            break
        }
        versionNum[i]="0"
        num, _ = strconv.Atoi(versionNum[i-1])
        num++
        versionNum[i-1]=fmt.Sprintf("%d",num)
        if num<9{
            break
        }
    }
    v:="v"
    for i :=0; i <= len(versionNum)-1 ;i++ {
        if i==0{
            v=v+versionNum[i]
            continue
        }
        v=v+"."+versionNum[i]
    }
    return v
}