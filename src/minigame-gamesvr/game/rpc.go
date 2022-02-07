package game

import (
	"log"

	"github.com/roada-go/roada"
)

type UserInsteadArgs struct {
}
type UserInsteadReply struct {
}

func UserInstead(rkey string) error {
	var args = UserInsteadArgs{}
	var reply = UserInsteadReply{}
	log.Printf("UserInstead, rkey=%s\n", rkey)
	err := roada.Call(rkey, "UserInstead", &args, &reply)
	if err != nil {
		//log.Printf("UserInstead failed, rkey=%s, error=%s\n", rkey, err.Error())
		return err
	}
	//log.Printf("UserInstead success, rkey=%s\n", rkey)
	return nil
}

func (agent *Agent) UserInstead(r *roada.Request, args *UserInsteadArgs, reply *UserInsteadReply) error {
	//now := time.Now().Unix()
	//log.Printf("instead instead instead instead")
	log.Printf("账号在其他地方登录 sessionid=%d\n", agent.session.ID())
	agent.kick("账号在其他地方登录")
	return nil
}
