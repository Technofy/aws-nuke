package cmd

type ItemPrinter interface {
	PrintItem(*Item)
	Flush()
}