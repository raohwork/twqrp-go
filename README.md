產生 TWQRP 內容的小工具

Since TWQRP is only used in Taiwan, the document is provided only in Traditional Chinese.

# WIP

此工具目前尚待開發與測試，API 介面隨時可能有不相容的變更，亦不建議在正式環境中使用

### 用法

```golang
NewTransfer("812", "12345678901234").Amount(150).Note("拆帳").String()
```

目前只支援產生轉帳內容，歡迎開 issue 或送 PR 提供更多 **經過驗證** 的內容

# License

MPL 2.0，詳見 LICENSE.txt
