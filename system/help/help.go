
package helper

func Menu(pushName string, prefix string) string {
	return `Hai, *` + pushName + `* 👋
*NB:* Bot ini masih tahap pengembangan!!

*- GROUP MENU -*
.linkgc
.resetlink
`
}

func BotAdm() string {
	return `*BOT NOT ADMIN*
	
Silahkan jadikan bot sebagai *admin* agar fitur ini dapat digunakan!!`
}

func Adm() string {
	return `*ADMIN ONLY*
	
Fitur ini hanya untuk admin !!`
}

func Own() string {
	return `*OWNER ONLY*
	
Fitur ini hanya untuk owner bot!!`
}

func Gc() string {
	return `*GROUP ONLY*
	
Fitur ini hanya untuk grup!!`
}
