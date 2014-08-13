package main

var (
	NUM map[int]string
)

func SetupNumerics() {
	NUM = make(map[int]string)
	NUM[RPL_WELCOME] = "Welcome to the Internet Relay Network %s!%s@%s"
	NUM[RPL_YOURHOST] = "Your host is %s, running %s version %s"
	NUM[RPL_CREATED] = "This server was created %s"
	NUM[RPL_MYINFO] = "%s %s %s %s"
	NUM[RPL_ENDOFWHO] = "%s :End of WHO list"
	NUM[RPL_LISTSTART] = "Channel :Users Name"
	NUM[RPL_LIST] = "%s %d :%s"
	NUM[RPL_LISTEND] = ":End of channel list"
	NUM[RPL_TOPIC] = "%s :%s"
	NUM[RPL_TOPICWHOTIME] = "%s %s %d"
	NUM[RPL_WHOREPLY] = "%s %s %s %s %s %s %s %s"
	NUM[RPL_NAMEPLY] = "= %s :%s"
	NUM[RPL_ENDOFNAMES] = "%s :End of /NAMES list."
	NUM[RPL_MOTD] = ":- %s"
	NUM[RPL_MOTDSTART] = ":- %s Message of the day -"
	NUM[RPL_ENDOFMOTD] = ":End of /MOTD command."
	NUM[ERR_NOSUCHCHANNEL] = "%s :No such channel"
	NUM[ERR_UNKNOWNCOMMAND] = "%s :Unkown command"
	NUM[ERR_NONICKNAMEGIVEN] = "No nickname given"
	NUM[ERR_NICKNAMEINUSE] = "%s :Nickname is already in use"
	NUM[ERR_ERRONEOUSNICKNAME] = "%s :Erroneous nickname"
	NUM[ERR_NOTREGISTERED] = ":You have not registered"
	NUM[ERR_NEEDMOREPARAMS] = "%s :Not enough parameters"
	NUM[ERR_ALREADYREGISTRED] = "You may not reregister"
}

const (
	RPL_WELCOME           = 001
	RPL_YOURHOST          = 002
	RPL_CREATED           = 003
	RPL_MYINFO            = 004
	RPL_ENDOFWHO          = 315
	RPL_LISTSTART         = 321
	RPL_LISTEND           = 323
	RPL_LIST              = 322
	RPL_TOPIC             = 332
	RPL_TOPICWHOTIME      = 333
	RPL_WHOREPLY          = 352
	RPL_NAMEPLY           = 353
	RPL_ENDOFNAMES        = 366
	RPL_MOTD              = 372
	RPL_MOTDSTART         = 375
	RPL_ENDOFMOTD         = 376
	ERR_NOSUCHCHANNEL     = 403
	ERR_UNKNOWNCOMMAND    = 421
	ERR_NONICKNAMEGIVEN   = 431
	ERR_ERRONEOUSNICKNAME = 432
	ERR_NICKNAMEINUSE     = 433
	ERR_NOTREGISTERED     = 451
	ERR_NEEDMOREPARAMS    = 461
	ERR_ALREADYREGISTRED  = 462
)
