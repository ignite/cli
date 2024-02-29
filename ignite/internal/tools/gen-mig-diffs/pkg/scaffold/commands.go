package scaffold

var defaultCommands = Commands{
	"chain": Command{
		Commands: []string{"chain example --no-module"},
	},
	"module": Command{
		Prerequisites: []string{"chain"},
		Commands:      []string{"module example --ibc"},
	},
	"list": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"list list1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"map": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"map map1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --index i1:string --module example --yes",
		},
	},
	"single": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"single single1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"type": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"type type1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"message": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"message message1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"query": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"query query1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints --module example --yes",
		},
	},
	"packet": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"packet packet1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --ack f1:string,f2:strings,f3:bool,f4:int,f5:ints,f6:uint,f7:uints,f8:coin,f9:coins --module example --yes",
		},
	},
}
