package color

type Color string

const (
	White Color = "\033[97m"
	Reset Color = "\033[0m"
	Red   Color = "\033[31m"
	Blue  Color = "\033[34m"
	Green Color = "\033[32m"
	//Purple Color = "\033[35m"
	//Cyan   Color = "\033[36m"
	//Gray   Color = "\033[37m"
	//Yellow Color = "\033[33m"
)
