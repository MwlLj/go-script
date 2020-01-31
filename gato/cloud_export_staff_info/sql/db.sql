#namespace db

/*
    @bref 获取所有的人员信息
    @is_brace true
    @in_isarr false
    @out_isarr true
    @out infrastructureUuid: string
    @out staffName: string
    @out address: string
    @out nation: string
    @out nationality: string
    @out nativePlace: string
    @out gender: string
*/
#define getAllStaffInfo
select sir.infrastructureUuid
, si.staffName, si.address, si.nation, si.nationality, si.nativePlace, si.gender
from t_vss_staff_infrastructure_rl as sir
inner join t_vss_staff_info as si
on sir.staffUuid = si.staffUuid;
#end

/*
    @bref 根据基建uuid获取基建信息
    @is_brace true
    @is_isarr false
    @out_isarr false
    @in infrastructureUuid: string
    @out infrastructureName: string
    @out parentUuid: string
*/
#define getInfrastructureInfoByUuid
select infrastructureName, parentUuid from t_vss_infrastructure_info where infrastructureUuid = {0};
#end

select infrastructureName, parentUuid from t_vss_infrastructure_info where infrastructureUuid = 'c52359de1b2d44728ec6fcd75d9c70b6';
