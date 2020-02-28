package main

const (
	Delete    = "Delete"
	FileName  = "FileName"
	Yes       = "Yes"
	No        = "No"
	ExportBtn = "Export"
	Extension = "Extension"
	Download  = "Download"
	Choose    = "Choose"
	BuildURL  = "BuildURL"
	URL       = "URL"
	SQL       = "SQL"
	Args      = "Args"
	Titles    = "Titles"
	Sheet     = "Sheet"
	AddSheet  = "Add Sheet"
	Untitled  = "Untitled"

	URLFormat = "%s:%s@tcp(%s:%s)/%s?charset=%s"
	URLRegex  = `(?m)\S+:\S+@tcp\(((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?):([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])\)/\S+\?[C|c]harset=\S+`
	SQLRegex  = `(?mi)select\s+(\*|(\S+(,?\s*\S+)*))\s+from\s+(\S+)(\s+where\s+.*)?`
)

var (
	charsets = []map[string]string{
		{"utf8": "[UTF-8 Unicode] - utf8_general_ci"},
		{"utf8mb4": "[UTF-8 Unicode] - utf8mb4_general_ci"},
		{"gbk": "[GBK Simplified Chinese] - gbk_chinese_ci"},
		{"ascii": "[US ASCII] - ascii_general_ci"},
		{"gb2312": "[GB2312 Simplified Chinese] - gb2312_chinese_ci"},
		{"greek": "[ISO 8859-7 Greek] - greek_general_ci"},
		{"utf16": "[UTF-16 Unicode] - utf16_general_ci"},
		{"utf16le": "[UTF-16LE Unicode] - utf16le_general_ci"},
		{"utf32": "[UTF-32 Unicode] - utf32_general_ci"},
		{"big5": "[Big5 Traditional Chinese] - big5_chinese_ci"},
		{"dec8": "[DEC West European] - dec8_swedish_ci"},
		{"cp850": "[DOS West European] - cp850_general_ci"},
		{"hp8": "[HP West European] - hp8_english_ci"},
		{"koi8r": "[KOI8-R Relcom Russian] - koi8r_general_ci"},
		{"latin1": "[cp1252 West European] - latin1_swedish_ci"},
		{"latin2": "[ISO 8859-2 Central European] - latin2_general_ci"},
		{"swe7": "[7bit Swedish] - swe7_swedish_ci"},
		{"ujis": "[EUC-JP Japanese] - ujis_japanese_ci"},
		{"sjis": "[Shift-JIS Japanese] - sjis_japanese_ci"},
		{"hebrew": "[ISO 8859-8 Hebrew] - hebrew_general_ci"},
		{"tis620": "[TIS620 Thai] - tis620_thai_ci"},
		{"euckr": "[EUC-KR Korean] - euckr_korean_ci"},
		{"koi8u": "[KOI8-U Ukrainian] - koi8u_general_ci"},
		{"cp1250": "[Windows Central European] - cp1250_general_ci"},
		{"latin5": "[ISO 8859-9 Turkish] - latin5_turkish_ci"},
		{"armscii8": "[ARMSCII-8 Armenian] - armscii8_general_ci"},
		{"ucs2": "[UCS-2 Unicode] - ucs2_general_ci"},
		{"cp866": "[DOS Russian] - cp866_general_ci"},
		{"keybcs2": "[DOS Kamenicky Czech-Slovak] - keybcs2_general_ci"},
		{"macce": "[Mac Central European] - macce_general_ci"},
		{"macroman": "[Mac West European] - macroman_general_ci"},
		{"cp852": "[DOS Central European] - cp852_general_ci"},
		{"latin7": "[ISO 8859-13 Baltic] - latin7_general_ci"},
		{"cp1251": "[Windows Cyrillic] - cp1251_general_ci"},
		{"cp1256": "[Windows Arabic] - cp1256_general_ci"},
		{"cp1257": "[Windows Baltic] - cp1257_general_ci"},
		{"binary": "[Binary pseudo charset] binary"},
		{"geostd8": "[GEOSTD8 Georgian] - geostd8_general_ci"},
		{"cp932": "[SJIS for Windows Japanese] - cp932_japanese_ci"},
		{"eucjpms": "[UJIS for Windows Japanese] - eucjpms_japanese_ci"},
		{"gb18030": "[China National Standard GB18030] - gb18030_chinese_ci"},
	}
)
