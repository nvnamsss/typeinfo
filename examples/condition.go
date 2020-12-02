package examples



//Condition comment
//
//Condition comment 2
type Condition struct {
	PoolID        int64   `gorm:"column:pool_id" json:"pool_id"`
	RewardID      int64   `gorm:"column:reward_id" json:"reward_id"`
	Name          string  `gorm:"column:name" json:"name"`
	IncludeReward []int64 `sql:"type:include_reward" json:"include_reward,omitempty"`
	ExcludeReward []int64 `sql:"type:exclude_reward" json:"exclude_reward,omitempty"`
}

//Get location comment
func (f *Condition) GetLocation() string {
	//implement here
	return "Ho Chi Minh"
}

func (f *Condition) ProcessReward() {

}

//IsLimitDeviceByDay comment
func (f *Condition) IsLimitDeviceByDay() bool {

	return false
}

func (f *Condition) IsLimitDeviceByCampaign() bool {
	return false
}

func (f *Condition) RewardClaimed(userID int64, campaignID int64) int64 {
	return 0
}

func (f *Condition) RewardClaimedTimes(userID int64, campaignID int64) int64 {
	return 0
}
