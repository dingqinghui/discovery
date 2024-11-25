/**
 * @Author: dingQingHui
 * @Description:
 * @File: member_list
 * @Version: 1.0.0
 * @Date: 2024/5/6 11:26
 */

package common

type Topology struct {
	EventId uint64
	Alive   []*Member
	Joined  []*Member
	Left    []*Member
}

type MemberList struct {
	Members     map[string]*Member
	LastEventId uint64
}

func NewMemberList() *MemberList {
	return &MemberList{
		Members: make(map[string]*Member),
	}
}

func (m *MemberList) UpdateClusterTopology(members map[string]*Member, lastEventId uint64) *Topology {
	if m.LastEventId >= lastEventId {
		return nil
	}
	tplg := &Topology{EventId: lastEventId}
	for _, member := range members {
		if _, ok := m.Members[member.GetID()]; ok {
			tplg.Alive = append(tplg.Alive, member)
		} else {
			tplg.Joined = append(tplg.Joined, member)
		}
	}
	for id := range m.Members {
		if _, ok := members[id]; !ok {
			tplg.Left = append(tplg.Left, m.Members[id])
		}
	}
	m.Members = members
	m.LastEventId = lastEventId
	return tplg
}
