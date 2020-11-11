## 這三個env是給程式裡面的獨立線程使用,key都是固定寫死,目前線程都關閉沒有在使用

- URI_IMEC
- URI_AMF
- URI_SON
- 這個設定完全和db裡面的config是分開獨立的, 不會互相影響

## db裡面的config設定只能有一筆,也就是說imec, son...等等的設定要整理後放在同一個[ ]裡面