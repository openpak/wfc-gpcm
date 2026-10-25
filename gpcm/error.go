package gpcm

import (
	"openpak/wfc-common/common"
	"openpak/wfc-common/gperror"
	"openpak/wfc-common/logging"
)

type GPError = gperror.GPError
type WWFCErrorMessage = gperror.WWFCErrorMessage

var MakeGPError = gperror.MakeGPError

const LangJapanese = gperror.LangJapanese
const LangEnglish = gperror.LangEnglish
const LangGerman = gperror.LangGerman
const LangFrench = gperror.LangFrench
const LangSpanish = gperror.LangSpanish
const LangItalian = gperror.LangItalian
const LangDutch = gperror.LangDutch
const LangSimpChinese = gperror.LangSimpChinese
const LangTradChinese = gperror.LangTradChinese
const LangKorean = gperror.LangKorean
const LangEnglishEU = gperror.LangEnglishEU
const LangFrenchEU = gperror.LangFrenchEU
const LangSpanishEU = gperror.LangSpanishEU

var ErrNone = gperror.ErrNone
var ErrGeneral = gperror.ErrGeneral
var ErrParse = gperror.ErrParse
var ErrNotLoggedIn = gperror.ErrNotLoggedIn
var ErrBadSessionKey = gperror.ErrBadSessionKey
var ErrDatabase = gperror.ErrDatabase
var ErrNetwork = gperror.ErrNetwork
var ErrForcedDisconnect = gperror.ErrForcedDisconnect
var ErrConnectionClosed = gperror.ErrConnectionClosed
var ErrUDPLayer = gperror.ErrUDPLayer
var ErrLogin = gperror.ErrLogin
var ErrLoginTimeout = gperror.ErrLoginTimeout
var ErrLoginBadNickname = gperror.ErrLoginBadNickname
var ErrLoginBadEmail = gperror.ErrLoginBadEmail
var ErrLoginBadPassword = gperror.ErrLoginBadPassword
var ErrLoginBadProfile = gperror.ErrLoginBadProfile
var ErrLoginProfileDeleted = gperror.ErrLoginProfileDeleted
var ErrLoginConnectionFailed = gperror.ErrLoginConnectionFailed
var ErrLoginServerAuthFailed = gperror.ErrLoginServerAuthFailed
var ErrLoginBadUniqueNick = gperror.ErrLoginBadUniqueNick
var ErrLoginBadPreAuth = gperror.ErrLoginBadPreAuth
var ErrLoginLoginTicketInvalid = gperror.ErrLoginLoginTicketInvalid
var ErrLoginLoginTicketExpired = gperror.ErrLoginLoginTicketExpired
var ErrNewUser = gperror.ErrNewUser
var ErrNewUserBadNickname = gperror.ErrNewUserBadNickname
var ErrNewUserBadPassword = gperror.ErrNewUserBadPassword
var ErrNewUserUniqueNickInvalid = gperror.ErrNewUserUniqueNickInvalid
var ErrNewUserUniqueNickInUse = gperror.ErrNewUserUniqueNickInUse
var ErrUpdateUserInfo = gperror.ErrUpdateUserInfo
var ErrUpdateUserInfoBadEmail = gperror.ErrUpdateUserInfoBadEmail
var ErrNewProfile = gperror.ErrNewProfile
var ErrNewProfileBadNickname = gperror.ErrNewProfileBadNickname
var ErrNewProfileBadOldNickname = gperror.ErrNewProfileBadOldNickname
var ErrUpdateProfile = gperror.ErrUpdateProfile
var ErrUpdateProfileBadNickname = gperror.ErrUpdateProfileBadNickname
var ErrAddFriend = gperror.ErrAddFriend
var ErrAddFriendBadFrom = gperror.ErrAddFriendBadFrom
var ErrAddFriendBadNew = gperror.ErrAddFriendBadNew
var ErrAddFriendAlreadyFriends = gperror.ErrAddFriendAlreadyFriends
var ErrAddFriendLocalBlock = gperror.ErrAddFriendLocalBlock
var ErrAddFriendBlocked = gperror.ErrAddFriendBlocked
var ErrAuthAdd = gperror.ErrAuthAdd
var ErrAuthAddBadFrom = gperror.ErrAuthAddBadFrom
var ErrAuthAddBadSignature = gperror.ErrAuthAddBadSignature
var ErrAuthAddLocalBlock = gperror.ErrAuthAddLocalBlock
var ErrAuthAddBlocked = gperror.ErrAuthAddBlocked
var ErrStatus = gperror.ErrStatus
var ErrMessage = gperror.ErrMessage
var ErrMessageNotFriends = gperror.ErrMessageNotFriends
var ErrMessageExtInfoNotSupported = gperror.ErrMessageExtInfoNotSupported
var ErrMessageFriendOffline = gperror.ErrMessageFriendOffline
var ErrGetProfile = gperror.ErrGetProfile
var ErrGetProfileBadProfile = gperror.ErrGetProfileBadProfile
var ErrDeleteFriend = gperror.ErrDeleteFriend
var ErrDeleteFriendNotFriends = gperror.ErrDeleteFriendNotFriends
var ErrDeleteProfile = gperror.ErrDeleteProfile
var ErrDeleteProfileLastProfile = gperror.ErrDeleteProfileLastProfile
var ErrSearch = gperror.ErrSearch
var ErrSearchConnectFailed = gperror.ErrSearchConnectFailed
var ErrSearchTimedOut = gperror.ErrSearchTimedOut
var ErrCheck = gperror.ErrCheck
var ErrCheckBadEmail = gperror.ErrCheckBadEmail
var ErrCheckBadNickname = gperror.ErrCheckBadNickname
var ErrCheckBadPassword = gperror.ErrCheckBadPassword
var ErrRevoke = gperror.ErrRevoke
var ErrRevokeNotFriends = gperror.ErrRevokeNotFriends
var ErrRegisterUniqueNick = gperror.ErrRegisterUniqueNick
var ErrRegisterUniqueNickTaken = gperror.ErrRegisterUniqueNickTaken
var ErrRegisterUniqueNickReserved = gperror.ErrRegisterUniqueNickReserved
var ErrRegisterUniqueNickBadNamespace = gperror.ErrRegisterUniqueNickBadNamespace
var ErrRegisterCDKey = gperror.ErrRegisterCDKey
var ErrRegisterCDKeyBadKey = gperror.ErrRegisterCDKeyBadKey
var ErrRegisterCDKeyAlreadySet = gperror.ErrRegisterCDKeyAlreadySet
var ErrRegisterCDKeyAlreadyTaken = gperror.ErrRegisterCDKeyAlreadyTaken
var ErrAddBlock = gperror.ErrAddBlock
var ErrAddBlockAlreadyBlocked = gperror.ErrAddBlockAlreadyBlocked
var ErrRemoveBlock = gperror.ErrRemoveBlock
var ErrRemoveBlockNotBlocked = gperror.ErrRemoveBlockNotBlocked
var WWFCMsgUnknownLoginError = gperror.WWFCMsgUnknownLoginError
var WWFCMsgDolphinSetupRequired = gperror.WWFCMsgDolphinSetupRequired
var WWFCMsgProfileBannedTOS = gperror.WWFCMsgProfileBannedTOS
var WWFCMsgProfileBannedTOSNow = gperror.WWFCMsgProfileBannedTOSNow
var WWFCMsgProfileRestricted = gperror.WWFCMsgProfileRestricted
var WWFCMsgProfileRestrictedNow = gperror.WWFCMsgProfileRestrictedNow
var WWFCMsgProfileRestrictedCustom = gperror.WWFCMsgProfileRestrictedCustom
var WWFCMsgProfileRestrictedNowCustom = gperror.WWFCMsgProfileRestrictedNowCustom
var WWFCMsgKickedGeneric = gperror.WWFCMsgKickedGeneric
var WWFCMsgKickedModerator = gperror.WWFCMsgKickedModerator
var WWFCMsgKickedRoomHost = gperror.WWFCMsgKickedRoomHost
var WWFCMsgKickedCustom = gperror.WWFCMsgKickedCustom
var WWFCMsgConsoleMismatch = gperror.WWFCMsgConsoleMismatch
var WWFCMsgConsoleMismatchDolphin = gperror.WWFCMsgConsoleMismatchDolphin
var WWFCMsgProfileIDInvalid = gperror.WWFCMsgProfileIDInvalid
var WWFCMsgProfileIDInUse = gperror.WWFCMsgProfileIDInUse
var WWFCMsgPayloadInvalid = gperror.WWFCMsgPayloadInvalid
var WWFCMsgInvalidELO = gperror.WWFCMsgInvalidELO

func (g *GameSpySession) replyError(gpErr GPError) {
	logging.Error(g.ModuleName, "Reply error:", gpErr.ErrorString)
	if !g.LoginInfoSet {
		msg := gpErr.GetMessage()
		// logging.Info(g.ModuleName, "Sending error message:", msg)
		err := common.SendPacket(ServerName, g.ConnIndex, []byte(msg))
		if gpErr.Fatal || err != nil {
			if err != nil {
				logging.Error("GPCM", "Failed to send packet:", err)
			}
			common.ShouldNotError(common.CloseConnection(ServerName, g.ConnIndex))
		}
		logging.Event("gpcm_returned_error", map[string]any{
			"profile_id":   g.User.ProfileId,
			"error_code":   gpErr.ErrorCode,
			"error_string": gpErr.ErrorString,
			"fatal":        gpErr.Fatal,
		})
		return
	}

	deviceId := g.User.RestrictedDeviceId
	if deviceId == 0 && len(g.User.NgDeviceId) > 0 {
		deviceId = g.User.NgDeviceId[0]
	}

	msg, wwfcErrorCode := gpErr.GetMessageTranslate(g.GameName, g.Region, g.Language, g.ConsoleFriendCode, deviceId)
	// logging.Info(g.ModuleName, "Sending error message:", msg)
	err := common.SendPacket(ServerName, g.ConnIndex, []byte(msg))
	if gpErr.Fatal || err != nil {
		if err != nil {
			logging.Error("GPCM", "Failed to send packet:", err)
		}
		common.ShouldNotError(common.CloseConnection(ServerName, g.ConnIndex))
	}

	logging.Event("gpcm_returned_error", map[string]any{
		"profile_id":         g.User.ProfileId,
		"error_code":         gpErr.ErrorCode,
		"error_string":       gpErr.ErrorString,
		"fatal":              gpErr.Fatal,
		"wiilink_error_code": wwfcErrorCode,
	})
}
