package core

import (
	"fmt"
)
type ChessMessage struct {
	Data []byte //每次对方移动位置的棋子的坐标
	/*ChesserType int64  //1表示红方，2表示黑方
	ChessType int64 //1=帅(将)，2=仕(士)，3=相(象),4=马,5=炮，6=车，7=兵(卒)
	hflag bool //表示红是否出棋
	bflag bool //表示黑是否出棋子
	Type int64 //1=求和，2=认输,3=继续对弈*/
	roomId string
	Conn *ChessConn
}

type Message struct {
	data []byte
	roomid string
	conn *Connection
}
type Hub struct {
	rooms map[string]map[*Connection]bool
	prepare chan Message
	unprepare chan Message
	broadcast   chan Message
	broadcastss chan Message
	warnings   chan Message
	login      chan Message
	unlogin chan Message
	kickoutroom chan Message
	warnmsg    chan Message
	games map[string]map[*ChessConn]bool
	register chan ChessMessage
	unregister chan ChessMessage
	illegalChess chan ChessMessage
	chessCoordinates chan ChessMessage


}

var H = Hub{
	broadcast:   make(chan Message),
	broadcastss: make(chan Message),
	warnings:    make(chan Message),
	warnmsg:     make(chan Message),
	login:       make(chan Message),
	unlogin:  	make(chan Message),
	kickoutroom: make(chan Message),
	rooms:       make(map[string]map[*Connection]bool),
	prepare:  make(chan Message),
	unprepare: make(chan Message),
	games: make(map[string]map[*ChessConn]bool),
	register: make(chan ChessMessage),
	unregister: make(chan ChessMessage),
	illegalChess: make(chan ChessMessage),
	chessCoordinates: make(chan ChessMessage),

}
func (H *Hub)Run() {
	for {
		select {
		case m := <-H.login: //传输链接
			conns := H.rooms[m.roomid]

			if conns == nil { //注册一个房间
				conn := make(map[*Connection]bool)
				H.rooms[m.roomid] = conn
				fmt.Println("在线人数==", len(conns))
				fmt.Println("room==", H.rooms)
			}

				H.rooms[m.roomid][m.conn] = true
				fmt.Println("在线人数:==", len(conns))
				fmt.Println("rooms:==", H.rooms)

			for con := range conns {
				delmsg := "系统消息：欢迎新伙伴加入" + m.roomid + "象棋室！！！"
				data := []byte(delmsg)
				select {
				case con.send <- data:
				}
			}
			/*case m:=<-H.prepare:
			prepares:=H.rooms[m.roomid]
			H.rooms[m.roomid][m.conn] = true
			fmt.Println("准备人数:==", len(prepares))
			fmt.Println("rooms:==", H.rooms)
			for prepare:=range prepares {
				delmsg := "系统消息: 对方已经准备" + m.roomid
				data := []byte(delmsg)
				select {
				case prepare.send <- data:
				}
			}*/
		/*case m:=<-H.unprepare:
			//取消准备
			prepares:=H.rooms[m.roomid]
			if prepares!=nil{
				if _,ok:=prepares[m.conn];ok{
					delete(prepares,m.conn)
					close(m.conn.send)
					for prepare:=range prepares{
						delemsg:="系统消息，对方取消了准备"+m.roomid
						data:=[]byte(delemsg)
						select {
						case prepare.send<-data:
						}
					}
				}
			}
*/


		case m := <-H.unlogin:

			conns := H.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn)
					close(m.conn.send)
					for conn := range conns {
						delmsg := "系统消息：有小伙伴离开了" + m.roomid + "象棋室！欢送！！！"
						data := []byte(delmsg)
						select {
						case conn.send <- data:
						}
						if len(conns) == 0 { // 链接都断开，删除房间
							delete(H.rooms, m.roomid)
						}

					}
				}
			}
		case m := <-H.kickoutroom:
			//三次不合法消息就踢出
			conns := H.rooms[m.roomid]
			notice := "由于您多次发送不合法信息,已被踢出！！！"
			select {
			case m.conn.send <- []byte(notice):
			}
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn)
					close(m.conn.send)
					if len(conns) == 0 {
						delete(H.rooms, m.roomid)
					}

				}

			}
			case m:=<-H.warnings: //不合法信息警告
				//发表警告
				conns:=H.rooms[m.roomid]
				if conns!=nil{
					if _,ok:=conns[m.conn];ok{
						notice := "警告:您发布不合法信息，将禁言5分钟，三次后将被踢出！！！"
						//starttime:=
						select {
						case m.conn.send <- []byte(notice):
						}
					}
				}
		case m:=<-H.warnmsg: //禁言中提示
			conns:=H.rooms[m.roomid]
			if conns!=nil{
				if _,ok:=conns[m.conn];ok{
					notice := "您还在禁言中,暂时不能发送信息！！！"
					select {
					case m.conn.send <- []byte(notice):
					}

				}
			}
		case m := <-H.broadcast:  //传输房间信息
			conns := H.rooms[m.roomid]
			for con := range conns {
				if con==m.conn{  //自己发送的信息，不用再发给自己
					continue
				}
				select {
				case con.send <- m.data:
				default:
					close(con.send)
					delete(conns, con)
					if len(conns) == 0 {
						delete(H.rooms, m.roomid)
					}
				}
			}
			case m:=<-H.register:
				conns:=H.games[m.roomId]
				if conns==nil{// 准备相应的房间
					conns=make(map[*ChessConn]bool)
					H.games[m.roomId]=conns
					fmt.Println("准备人数:==",len(conns))
					fmt.Println("rooms:==",H.rooms)
				}
				H.games[m.roomId][m.Conn]=true
				fmt.Println("准备人数:==",len(conns))
				fmt.Println("rooms:==",H.rooms)
				for conn:=range conns{
					dlemsg:="系统消息: 对方已经准备"
					data:=[]byte(dlemsg)
					select {
					case conn.send<-data:
					}

				}
				case m:=<-H.unregister: //不准备
				conns:=H.games[m.roomId]
				if conns!=nil{
					if _,ok:=conns[m.Conn];ok{
						delete(conns,m.Conn)
						close(m.Conn.send)
						for conn:=range conns{
							dlemsg:="系统消息：对面没有准备"
							data:=[]byte(dlemsg)
							select {
							case conn.send<-data:
							}
						}
						if len(conns)==0{//关闭游戏
							delete(H.games,m.roomId)
						}

					}
				}
		case m:=<-H.illegalChess:
			conns:=H.games[m.roomId]//拿到链接
			if conns!=nil{
				if _,ok:=conns[m.Conn];ok{
					notice := "现在是对方的回合,暂时不能出棋！！！"
					select {
					case m.Conn.send <- []byte(notice):
					}

				}
			}

		}
	}
}
