create function get_avg_duration(_machine_id int,_mitarbeiter_id int
								,_auftrags_id int,_limit int)
returns float
language plpgsql
as
$$
declare
   avg_duration integer;
begin

 if _limit <0 then
   SELECT avg(diff) into avg_duration FROM
   (SELECT diff,id FROM counters where machine_id=_machine_id and mitarbeiter_id=_mitarbeiter_id and auftrags_id=_auftrags_id ORDER BY id DESC LIMIT _limit*-1)as Laufzeit;
   return avg_duration;
   elsif _limit>0 then
   SELECT avg(diff) into avg_duration FROM
   (SELECT diff,id FROM counters where machine_id=_machine_id and mitarbeiter_id=_mitarbeiter_id and auftrags_id=_auftrags_id ORDER BY id asc LIMIT _limit)as Laufzeit;
   return avg_duration;
   elsif _limit = 0 then
   SELECT avg(diff) into avg_duration FROM
   (SELECT diff,id FROM counters where machine_id=_machine_id and mitarbeiter_id=_mitarbeiter_id and auftrags_id=_auftrags_id ORDER BY id asc)as Laufzeit;
   return avg_duration;
   END if;

end;
$$;
