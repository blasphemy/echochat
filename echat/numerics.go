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
	NUM[RPL_ISUPPORT] = "%s :are supported by this server"
	NUM[RPL_LUSERCLIENT] = ":There are %d users and %d services on %d servers"
	NUM[RPL_LUSEROP] = "%d :operator(s) online"
	NUM[RPL_LUSERUNKNOWN] = "%d :unknown connection(s)"
	NUM[RPL_LUSERCHANNELS] = "%d :channels formed"
	NUM[RPL_LUSERME] = ":I have %d clients and %d servers"
	NUM[RPL_ENDOFWHO] = "%s :End of WHO list"
	NUM[RPL_LISTSTART] = "Channels :Users Name"
	NUM[RPL_LIST] = "%s %d :%s"
	NUM[RPL_NOTOPIC] = "%s :No topic is set"
	NUM[RPL_LISTEND] = ":End of channel list"
	NUM[RPL_CHANNELMODEIS] = "%s +%s"
	NUM[RPL_CREATIONTIME] = "%s %d"
	NUM[RPL_TOPIC] = "%s :%s"
	NUM[RPL_TOPICWHOTIME] = "%s %s %d"
	NUM[RPL_WHOREPLY] = "%s %s %s %s %s %s %s %s"
	NUM[RPL_NAMEPLY] = "= %s :%s"
	NUM[RPL_ENDOFNAMES] = "%s :End of /NAMES list."
	NUM[RPL_MOTD] = ":- %s"
	NUM[RPL_MOTDSTART] = ":- %s Message of the day -"
	NUM[RPL_ENDOFMOTD] = ":End of /MOTD command."
	NUM[ERR_NOSUCHCHANNEL] = "%s :No such channel"
	NUM[ERR_NOSUCHNICK] = "%s :No such nick"
	NUM[ERR_CANNOTSENDTOCHAN] = "%s: Cannot send to channel."
	NUM[ERR_UNKNOWNCOMMAND] = "%s :Unkown command"
	NUM[ERR_NONICKNAMEGIVEN] = "No nickname given"
	NUM[ERR_NICKNAMEINUSE] = "%s :Nickname is already in use"
	NUM[ERR_USERNOTINCHANNEL] = "%s %s :isn't on that channel"
	NUM[ERR_ERRONEOUSNICKNAME] = "%s :Erroneous nickname"
	NUM[ERR_NOTREGISTERED] = ":You have not registered"
	NUM[ERR_NEEDMOREPARAMS] = "%s :Not enough parameters"
	NUM[ERR_ALREADYREGISTRED] = "You may not reregister"
	NUM[ERR_CHANOPRIVSNEEDED] = "%s :You do not have the required status to perform this action"
	NUM[RPL_YOUREOPER] = ":You are now an IRC operator"
}

const (
	RPL_WELCOME           = 001
	RPL_YOURHOST          = 002
	RPL_CREATED           = 003
	RPL_MYINFO            = 004
	RPL_ISUPPORT          = 005
	RPL_LUSERCLIENT       = 251
	RPL_LUSEROP           = 252
	RPL_LUSERUNKNOWN      = 253
	RPL_LUSERCHANNELS     = 254
	RPL_LUSERME           = 255
	RPL_ENDOFWHO          = 315
	RPL_LISTSTART         = 321
	RPL_LIST              = 322
	RPL_LISTEND           = 323
	RPL_CHANNELMODEIS     = 324
	RPL_CREATIONTIME      = 329
	RPL_NOTOPIC           = 331
	RPL_TOPIC             = 332
	RPL_TOPICWHOTIME      = 333
	RPL_WHOREPLY          = 352
	RPL_NAMEPLY           = 353
	RPL_ENDOFNAMES        = 366
	RPL_MOTD              = 372
	RPL_MOTDSTART         = 375
	RPL_ENDOFMOTD         = 376
	ERR_NOSUCHCHANNEL     = 403
	ERR_CANNOTSENDTOCHAN  = 404
	ERR_UNKNOWNCOMMAND    = 421
	ERR_NONICKNAMEGIVEN   = 431
	ERR_ERRONEOUSNICKNAME = 432
	ERR_NICKNAMEINUSE     = 433
	ERR_USERNOTINCHANNEL  = 441
	ERR_NOTREGISTERED     = 451
	ERR_NEEDMOREPARAMS    = 461
	ERR_ALREADYREGISTRED  = 462
	ERR_CHANOPRIVSNEEDED  = 482
	ERR_NOSUCHNICK        = 401
	RPL_YOUREOPER         = 381
)
