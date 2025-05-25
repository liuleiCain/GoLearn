package switch_case

func SwitchCase() {
	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			// fallthrough将会执行case 0和case 1
			println("case 0---", i)
			fallthrough
		case 1:
			println("case 1---", i)
			break
		case 2:
			// 默认会break
			println("case 2---", i)
		case 3:
			// fallthrough将会执行case 3和case default
			println("case 3---", i)
			fallthrough
		default:
			println("case default---", i)
		}
	}
}

func SwitchCaseMulti() {
	for i := 0; i < 4; i++ {
		switch i {
		case 0, 1:
			println("case 0, 1---", i)
			fallthrough
		case 2, 3:
			println("case 2, 3---", i)
		default:
			println("case default---", i)
		}
	}
}
