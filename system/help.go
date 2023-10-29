/*
###################################
# Name: Mywa BOT                  #
# Version: 1.0.1                  #
# Developer: Amirul Dev           #
# Library: waSocket               #
# Contact: 085157489446           #
###################################
# Thanks to: 
# Vnia
*/
package system

func Menu(pushName string, prefix string) string {
    return `Hai, *` + pushName + `* 👋
*NB:* Bot ini masih tahap pengembangan!!

𖥔 Database: Firebase
𖥔 Library: MywaGO
𖥔 Language: GoLang
𖥔 Size Script: 1.8M

⌬  MAIN MENU ⌬ 
⦿ .menu

⌬ GROUP MENU ⌬ 
⦿ .linkgc
⦿ .resetlink

⌬  TOOLS MENU ⌬ 
⦿ .ai
⦿ .bard

⌬ OWNER MENU ⌬ 
⦿ $
⦿ .restart
⦿ .backup

2023 © 𝑳𝒊𝒈𝒉𝒕𝒘𝒆𝒊𝒈𝒉𝒕 𝒘𝒉𝒂𝒕𝒔𝒂𝒑𝒑 𝒃𝒐𝒕
𝒎𝒂𝒅𝒆 𝒃𝒚 𝑨𝒎𝒊𝒓𝒖𝒍 𝑫𝒆𝒗 ×͜×
- www.amirull.dev`
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