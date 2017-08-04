package main

import (
	"github.com/golang/protobuf/proto"
	"share/pb"
)

type ProtoMessageFunc func() proto.Message

var ProtocolHandlers = make(map[pb.Cmdid]ProtoMessageFunc)

func init() {
	add(pb.Cmdid_CmdLogin, func() proto.Message { return new(pb.LoginReq) })
	add(pb.Cmdid_CmdHeartbeat, func() proto.Message { return new(pb.Heartbeat) })
	add(pb.Cmdid_CmdChangeHeartbeatType, func() proto.Message { return new(pb.ChangeHeartbeatType) })

	// game
	add(pb.Cmdid_CmdAssignRole, func() proto.Message { return &pb.AssignRole{} })
	add(pb.Cmdid_CmdProtect, func() proto.Message { return &pb.Protect{} })
	add(pb.Cmdid_CmdKillConfirm, func() proto.Message { return &pb.KillConfirm{} })
	add(pb.Cmdid_CmdCheck, func() proto.Message { return &pb.Check{} })
	add(pb.Cmdid_CmdCure, func() proto.Message { return &pb.Cure{} })
	add(pb.Cmdid_CmdPoison, func() proto.Message { return &pb.Poison{} })
	add(pb.Cmdid_CmdStopSpeech, func() proto.Message { return &pb.StopSpeech{} })
	add(pb.Cmdid_CmdCaptainVote, func() proto.Message { return &pb.CaptainVote{} })
	add(pb.Cmdid_CmdCaptainVoteResult, func() proto.Message { return &pb.CaptainVoteResult{} })
	add(pb.Cmdid_CmdGodMessage, func() proto.Message { return &pb.GodMessage{} })
	add(pb.Cmdid_CmdHunt, func() proto.Message { return &pb.Hunt{} })
	add(pb.Cmdid_CmdChooseSuccessor, func() proto.Message { return &pb.ChooseCaptainSuccessor{} })
	add(pb.Cmdid_CmdGiveUpElection, func() proto.Message { return &pb.GiveUpElection{} })
	add(pb.Cmdid_CmdDecideSpeechOrder, func() proto.Message { return &pb.DecideSpeechOrder{} })
	add(pb.Cmdid_CmdShowEliminateVoteResult, func() proto.Message { return &pb.ShowEliminateVoteResult{} })
	add(pb.Cmdid_CmdNightTime, func() proto.Message { return &pb.NightTime{} })
	add(pb.Cmdid_CmdDayTime, func() proto.Message { return &pb.DayTime{} })
	add(pb.Cmdid_CmdShowCaptainCandidates, func() proto.Message { return &pb.ShowCaptainCandidates{} })
	add(pb.Cmdid_CmdSpeech, func() proto.Message { return &pb.Speech{} })
	add(pb.Cmdid_CmdEliminateVote, func() proto.Message { return &pb.EliminateVote{} })
	add(pb.Cmdid_CmdPlayerOffline, func() proto.Message { return &pb.PlayerOffline{} })
	add(pb.Cmdid_CmdPlayerOnline, func() proto.Message { return &pb.PlayerOnline{} })
	add(pb.Cmdid_CmdWerewolfSuicide, func() proto.Message { return &pb.WerewolfSuicide{} })
	add(pb.Cmdid_CmdGameOver, func() proto.Message { return &pb.GameOver{} })
	add(pb.Cmdid_CmdKillPrepare, func() proto.Message { return &pb.KillPrepare{} })
	add(pb.Cmdid_CmdNewCaptain, func() proto.Message { return &pb.NewCaptain{} })
	add(pb.Cmdid_CmdStopAllSpeech, func() proto.Message { return &pb.StopAllSpeech{} })
	add(pb.Cmdid_CmdKillResult, func() proto.Message { return &pb.KillResult{} })
	add(pb.Cmdid_CmdCheckResult, func() proto.Message { return &pb.CheckResult{} })
	add(pb.Cmdid_CmdHuntResult, func() proto.Message { return &pb.HuntResult{} })
	add(pb.Cmdid_CmdMuteMic, func() proto.Message { return &pb.MuteMic{} })
	add(pb.Cmdid_CmdUnmuteMic, func() proto.Message { return &pb.UnmuteMic{} })
	add(pb.Cmdid_CmdAlertMessage, func() proto.Message { return &pb.AlertMessage{} })
	add(pb.Cmdid_CmdOperMessage, func() proto.Message { return &pb.OperMessage{} })
	add(pb.Cmdid_CmdPlayerEscape, func() proto.Message { return &pb.PlayerEscape{} })
	add(pb.Cmdid_CmdCurrentGameStatus, func() proto.Message { return &pb.CurrentGameStatus{} })
	add(pb.Cmdid_CmdGrabRole, func() proto.Message { return &pb.GrabRole{} })
	add(pb.Cmdid_CmdUseOvertimeCard, func() proto.Message { return &pb.UseOvertimeCard{} })
	add(pb.Cmdid_CmdSomeoneGrabRole, func() proto.Message { return &pb.SomeoneGrabRole{} })
	add(pb.Cmdid_CmdUseOvertimeCardResult, func() proto.Message { return &pb.UseOvertimeCardResult{} })
	add(pb.Cmdid_CmdSpeecherPKStart, func() proto.Message { return &pb.SpeecherPKStart{} })
	add(pb.Cmdid_CmdPlayerGoDie, func() proto.Message { return &pb.PlayerGoDie{} })
	add(pb.Cmdid_CmdEmptySpeak, func() proto.Message { return &pb.EmptySpeak{} })
	add(pb.Cmdid_CmdPenguinStartFreeze, func() proto.Message { return &pb.PenguinStartFreeze{} })
	add(pb.Cmdid_CmdPenguinSelfFreezeResult, func() proto.Message { return &pb.PenguinSelfFreezeResult{} })
	add(pb.Cmdid_CmdPenguinFreezeResult, func() proto.Message { return &pb.PenguinFreezeResult{} })
	add(pb.Cmdid_CmdCrowStartGossip, func() proto.Message { return &pb.CrowStartGossip{} })
	add(pb.Cmdid_CmdCrowSelfGossipResult, func() proto.Message { return &pb.CrowSelfGossipResult{} })
	add(pb.Cmdid_CmdBearShoutResult, func() proto.Message { return &pb.BearShoutResult{} })
	add(pb.Cmdid_CmdCrowGossipResult, func() proto.Message { return &pb.CrowGossipResult{} })
	add(pb.Cmdid_CmdBlackWolfKingStartHunt, func() proto.Message { return &pb.BlackWolfKingStartHunt{} })
	add(pb.Cmdid_CmdBlackWolfKingHuntResult, func() proto.Message { return &pb.BlackWolfKingHuntResult{} })

	// room
	add(pb.Cmdid_CmdPlayerSitDown, func() proto.Message { return &pb.PlayerSitDown{} })
	add(pb.Cmdid_CmdPlayerWitness, func() proto.Message { return &pb.PlayerWitness{} })
	add(pb.Cmdid_CmdPlayerPrepared, func() proto.Message { return &pb.PlayerPrepared{} })
	add(pb.Cmdid_CmdPlayerCanceled, func() proto.Message { return &pb.PlayerCanceled{} })
	add(pb.Cmdid_CmdPlayerLeave, func() proto.Message { return &pb.PlayerLeave{} })
	add(pb.Cmdid_CmdRoomMessage, func() proto.Message { return &pb.RoomMessage{} })
	add(pb.Cmdid_CmdRoomOwnerChanged, func() proto.Message { return &pb.RoomOwnerChanged{} })
	add(pb.Cmdid_CmdPlayerJoin, func() proto.Message { return &pb.PlayerJoin{} })
	add(pb.Cmdid_CmdMatchRoom, func() proto.Message { return &pb.MatchRoom{} })
	add(pb.Cmdid_CmdGameStart, func() proto.Message { return &pb.GameStart{} })
	add(pb.Cmdid_CmdRoomSettingChanged, func() proto.Message { return &pb.RoomSettingChanged{} })
	add(pb.Cmdid_CmdRoomKick, func() proto.Message { return &pb.RoomKick{} })
	add(pb.Cmdid_CmdSpeechStart, func() proto.Message { return &pb.SpeechStart{} })
	add(pb.Cmdid_CmdSpeechStop, func() proto.Message { return &pb.SpeechStop{} })
	//add(pb.Cmdid_CmdChatOutside, func() proto.Message { return &pb.ChatOutside{} })
	add(pb.Cmdid_CmdConfirmHangup, func() proto.Message { return &pb.ConfirmHangup{} })
	add(pb.Cmdid_CmdRoomInvitation, func() proto.Message { return &pb.RoomInvitation{} })
	//add(pb.Cmdid_CmdChatOutsideNew, func() proto.Message { return &pb.ChatOutsideNew{} })
	//add(pb.Cmdid_CmdChatDeadPlayer, func() proto.Message { return &pb.ChatDeadPlayer{} })
}

func add(cmdId pb.Cmdid, fn ProtoMessageFunc) {
	ProtocolHandlers[cmdId] = fn
}
