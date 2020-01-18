package main

import (
    "main/handler"
)

func main() {
    obj := handler.NewByParagraph("resource/test.xlsx", "resource/test.xlsx")
    obj.Calc()
}
