/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/8/31 18:04
 * @Desc: 私聊消息实体
 */

package group

import (
	"errors"
	
	"github.com/webzh/tencent-im/internal/entity"
)

var (
	errNotSetSender   = errors.New("message's sender not set")
	errNotSetSendTime = errors.New("message's send time not set")
)

type (
	// MsgOnlineOnlyFlag 只发送在线成员标识
	MsgOnlineOnlyFlag int
	
	// MsgPriority 消息优先级
	MsgPriority string
	
	// MsgStatus 消息状态
	MsgStatus int
)

const (
	MsgOnlineOnlyFlagNo  MsgOnlineOnlyFlag = 0 // 发送所有成员
	MsgOnlineOnlyFlagYes MsgOnlineOnlyFlag = 1 // 仅发送在线成员
	
	MsgPriorityHigh   MsgPriority = "High"   // 高优先级消息
	MsgPriorityNormal MsgPriority = "Normal" // 普通优先级消息
	MsgPriorityLow    MsgPriority = "Low"    // 低优先级消息
	MsgPriorityLowest MsgPriority = "Lowest" // 最低优先级消息
	
	MsgStatusNormal  MsgStatus = 0 // 正常消息
	MsgStatusInvalid MsgStatus = 1 // 被删除或者消息过期的消息
	MsgStatusRevoked MsgStatus = 2 // 被撤回的消息
	
	AtAllMembersFlag = "@all" // @所有成员的标识
)

type Message struct {
	entity.Message
	priority         MsgPriority       // 消息的优先级
	onlineOnlyFlag   MsgOnlineOnlyFlag // 仅发送在线成员标识
	sendTime         int64             // 消息发送时间
	timestamp        int64             // 消息时间戳，UNIX 时间戳（单位：秒）
	seq              int               // 消息序列号
	status           MsgStatus         // 消息状态
	customData       interface{}       // 自定义数据
	sendControls     map[string]bool   // 发送消息控制
	callbackControls map[string]bool   // 禁用回调
	atMembers        map[string]bool   // @用户
}

func NewMessage() *Message {
	return &Message{}
}

// SetPriority 设置消息优先级
func (m *Message) SetPriority(priority MsgPriority) {
	m.priority = priority
}

// GetPriority 获取消息优先级
func (m *Message) GetPriority() MsgPriority {
	return m.priority
}

// SetCustomData 设置自定义数据
func (m *Message) SetCustomData(data interface{}) {
	m.customData = data
}

// GetCustomData 获取自定义数据
func (m *Message) GetCustomData() interface{} {
	return m.customData
}

// SetOnlineOnlyFlag 设置仅发送在线成员标识
func (m *Message) SetOnlineOnlyFlag(flag MsgOnlineOnlyFlag) {
	m.onlineOnlyFlag = flag
}

// GetOnlineOnlyFlag 获取仅发送在线成员标识
func (m *Message) GetOnlineOnlyFlag() MsgOnlineOnlyFlag {
	return m.onlineOnlyFlag
}

// SetSendTime 设置发送时间
func (m *Message) SetSendTime(sendTime int64) {
	m.sendTime = sendTime
}

// GetSendTime 获取发送时间
func (m *Message) GetSendTime() int64 {
	return m.sendTime
}

// GetStatus 获取消息状态
func (m *Message) GetStatus() MsgStatus {
	return m.status
}

// SetForbidBeforeSendMsgCallback 设置禁止发消息前回调
func (m *Message) SetForbidBeforeSendMsgCallback() {
	if m.callbackControls == nil {
		m.callbackControls = make(map[string]bool, 0)
	}
	m.callbackControls["ForbidBeforeSendMsgCallback"] = true
}

// SetForbidAfterSendMsgCallback 设置禁止发消息后回调
func (m *Message) SetForbidAfterSendMsgCallback() {
	if m.callbackControls == nil {
		m.callbackControls = make(map[string]bool, 0)
	}
	m.callbackControls["ForbidAfterSendMsgCallback"] = true
}

// GetForbidCallbackControl 获取消息回调禁止开关
func (m *Message) GetForbidCallbackControl() (controls []string) {
	if m.callbackControls != nil {
		if n := len(m.callbackControls); n > 0 {
			controls = make([]string, 0, n)
			for k := range m.callbackControls {
				controls = append(controls, k)
			}
		}
	}
	
	return
}

// SetNoUnread 设置该条消息不计入未读数
func (m *Message) SetNoUnread() {
	if m.sendControls == nil {
		m.sendControls = make(map[string]bool, 0)
	}
	m.sendControls["NoUnread"] = true
}

// SetNoLastMsg 设置该条消息不更新会话列表
func (m *Message) SetNoLastMsg() {
	if m.sendControls == nil {
		m.sendControls = make(map[string]bool, 0)
	}
	m.sendControls["NoLastMsg"] = true
}

// GetSendMsgControl 获取消息发送控制选项
func (m *Message) GetSendMsgControl() (controls []string) {
	if m.sendControls != nil {
		if n := len(m.sendControls); n > 0 {
			controls = make([]string, 0, n)
			for k := range m.sendControls {
				controls = append(controls, k)
			}
		}
	}
	
	return
}

// AtMembers @某个成员
func (m *Message) AtMembers(userId ...string) {
	if m.atMembers == nil {
		m.atMembers = make(map[string]bool)
	}
	
	for _, id := range userId {
		m.atMembers[id] = true
	}
}

// AtAllMembers @所有成员
func (m *Message) AtAllMembers() {
	m.AtMembers(AtAllMembersFlag)
}

// ClearAtMembers 清空所有的的@成员
func (m *Message) ClearAtMembers() {
	m.atMembers = nil
}

// GetTimestamp 获取消息的时间戳
func (m *Message) GetTimestamp() int64 {
	return m.timestamp
}

// 检测发送错误
func (m *Message) checkSendError() (err error) {
	if err = m.CheckBodyArgError(); err != nil {
		return
	}
	
	return
}

// 检测导入错误
func (m *Message) checkImportError() (err error) {
	if m.GetSender() == "" {
		return errNotSetSender
	}
	
	if m.sendTime == 0 {
		return errNotSetSendTime
	}
	
	if err = m.CheckBodyArgError(); err != nil {
		return
	}
	
	return
}
