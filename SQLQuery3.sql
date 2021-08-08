select distinct(v.id) from VerifierUser v JOIN

 

(select r.UserId from ReportRequest r JOIN ReportRequest r1 on r.UserId=r1.UserId and DATEDIFF(month,r.CreatedOn,getdate())<6 group by r.UserId ) re
on 
(v.id<>re.UserId and V.VerifierType in (404,403) and v.IsActive=1) order by v.id

