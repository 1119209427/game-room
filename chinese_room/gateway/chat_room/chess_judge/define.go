package chess_judge

// FlagArr 判断是否有棋子，如果有就变为1.没有就是0
var FlagArr [10][9]int
//红方棋子的位置储存
var (
	RRooks [][]int //存储车的位置
	RKnight [][]int//存储马的位置
	RElephants [][]int//存在象的位置
	RMandarins [][]int //存储士的位置
	RKing [][]int //存在将的位置
	RCannons [][]int //存储炮的位置
	RPawns [][]int //存储兵的位置

)
//黑方棋子的位置存储
var (
	BRooks [][]int //存储车的位置
	BKnight [][]int//存储马的位置
	BElephants [][]int//存在象的位置
	BMandarins [][]int //存储士的位置
	BKing [][]int //存在将的位置
	BCannons [][]int //存储炮的位置
	BPawns [][]int //存储
)

// InitChess 初始化棋盘
func InitChess(){
	//9列，10行
	//为棋盘的第一行赋值
	for j:=0;j<9;j++{
			FlagArr[0][j]=1
		}
	for j:=0;j<9;j++{
		FlagArr[8][j]=1
	}
	//为炮所在的位置赋值
	FlagArr[2][1]=1
	FlagArr[2][7]=1
	RCannons=[][]int{[]int{2,1},[]int{2,7}}
	FlagArr[7][1]=1
	FlagArr[7][7]=1
	BCannons=[][]int{[]int{7,1},[]int{7,7}}
	//为兵所在的位置赋值
	for i:=0;i<9;i=i+2{
		FlagArr[3][i]=1
		RPawns=append(RPawns,[]int{3,i})
		FlagArr[6][i]=1
		BPawns=append(BPawns,[]int{6,i})
	}
	//为车赋值
	for i:=0;i<9;i+=8{
		RRooks=append(RRooks,[]int{0,i})
		BPawns=append(BPawns,[]int{9,i})
	}
	//为马赋值
	for i:=1;i<9;i+=6{
		RKnight=append(RKnight,[]int{0,i})
		BKnight=append(BKnight,[]int{9,i})
	}
	//把象的位置存储
	for i:=2;i<9;i+=4{
		RElephants=append(RElephants,[]int{0,i})
		BElephants=append(BElephants,[]int{9,i})
	}
	//存储士的位置
	for i:=3;i<9;i+=2{
		RMandarins=append(RMandarins,[]int{0,i})
		BMandarins=append(BMandarins,[]int{9,i})
	}
	RKing=[][]int{[]int{0,4}}
	BKing=[][]int{[]int{9,4}}
}


