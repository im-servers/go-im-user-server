
—
user_server.protouser_server"
GetUsersReq
id (Rid"8
GetUsersReply'
users (2.user_server.UserRusers"‚
User
id (Rid
name (	Rname
number (	Rnumber
gender (	Rgender
phone (	Rphone
age (Rage"
IdsReq
ids (Rids"³
CreateUserReq
name (	Rname
number (	Rnumber
gender (	Rgender
phone (	Rphone
age (Rage
password (	Rpassword
deviceID (	RdeviceID"K
CreateUserReply
id (Rid(
token (2.user_server.TokenRtoken"X
loginReq
phone (	Rphone
password (	Rpassword
deviceID (	RdeviceID"F

loginReply
id (Rid(
token (2.user_server.TokenRtoken"q
Token 
accessToken (	RaccessToken"
accessExpire (RaccessExpire"
refreshAfter (RrefreshAfter"a
authUserReq
id (Rid 
accessToken (	RaccessToken 
deviceToken (	RdeviceToken"S
authUserReply 
accessToken (	RaccessToken 
deviceToken (	RdeviceToken".
getDeviceTokensByUserIDReq
ids (Rids"î
getDeviceTokensByUserIDReplyk
userDeviceTokens (2?.user_server.getDeviceTokensByUserIDReply.UserDeviceTokensEntryRuserDeviceTokensa
UserDeviceTokensEntry
key (Rkey2
value (2.user_server.ListDeviceTokenRvalue:8")
ListDeviceToken
values (	Rvalues"-
GetUsersByTokensReq
tokens (	Rtokens"@
GetUsersByTokensReply'
users (2.user_server.UserRusers2Ô
user@
getUsers.user_server.GetUsersReq.user_server.GetUsersReplyX
getUsersByTokens .user_server.GetUsersByTokensReq".user_server.GetUsersByTokensReplyF

createUser.user_server.CreateUserReq.user_server.CreateUserReply7
login.user_server.loginReq.user_server.loginReply@
authUser.user_server.authUserReq.user_server.authUserReplym
getDeviceTokensByUserID'.user_server.getDeviceTokensByUserIDReq).user_server.getDeviceTokensByUserIDReplyBZ./user_serverbproto3