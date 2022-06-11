package chess_judge
func RooksMove(start,end string,flag int64)bool{//车的移动
	//車”的走棋规则是只能沿着横向或纵向走直线、且不能越过其他棋子。
	//先判断车的坐标开始坐标是否合法
	/*if flag==1{

		for _,v:=range RRooks{
			for i:=0;i<2;i++{
				v[i]=start[i]
			}
		}

	}*/

	return true

}
func KingMove(start,end string,flag int64)bool{
	return true
}
func MandarinsMove(start,end string,flag int64)bool{
	return true
}
func PawnsMove(start,end string,flag int64)bool{
	return true
}
func CannonsMove(start,end string,flag int64)bool{
	return true
}
func ElephantsMove(start,end string,flag int64)bool{
	return true
}