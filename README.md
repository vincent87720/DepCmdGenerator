# Classroom Deployment Command Generator
由於還原卡的故障導致IP位址和電腦名稱無法自動設定，藉由此程式可以生成設定批次檔，減少設定步驟及時間

## 環境
- Windows10 x64 1909

## 設定
本程式藉由`settings.json`檔案設定指令內容，說明如下
```json
{
    "StartHostInfo":"H705-01",
    "StartIPAddr":"192.168.0.1",
    "Mask":"255.255.255.0",
    "Gateway":"192.168.0.254",
    "DNSServer":[
        "192.168.1.66",
        "192.168.1.77"
    ],
    "Qty":50,
    "Interface":"乙太網路"
}
```
`StartHostInfo`
填入第一台設備的電腦名稱，'-'後面的數字會依照`Qty`數量多寡遞增

`StartIPAddr`
填入第一台設備的IP位址，最後一個'.'後面的數字會依照`Qty`數量多寡遞增

`Mask`
所有設備的子網路遮罩

`Gateway`
所有設備的閘道器(網關)位址

`DNSServer`
所有設備的DNSServer位址

`Qty`
要生產的批次檔數量

`Interface`
要進行設定的網路介面卡名稱


## 執行

*注意! settings.json與DepCmdGenerator.exe必須位於相同目錄*

### 步驟一
點擊`DepCmdGenerator.exe`執行可執行檔

### 步驟二
將生成出來的批次檔(.bat)移動到要設定的目標電腦，並對檔案點擊右鍵，使用***以系統管理員身分執行***選項啟動設定批次檔

## 驗證

開啟命令提示字元(cmd)程式，輸入`ipconfig/all`

### 電腦名稱設定
檢查`Windows IP 設定`項目的`主機名稱`是否已被修改

### 網路設定
查看欲修改的網卡項目是否已被成功修改
