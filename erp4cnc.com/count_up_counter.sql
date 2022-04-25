drop function count_up_counter;
create function count_up_counter(_machine_id int,_mitarbeiter_id int
								,_auftrags_id int)
returns int
language plpgsql
as
$$
declare
   _diff integer;
   _count integer;
   _id integer;
   last_entry timestamp;
   new_entry  timestamp;
begin
   select created_at,count into last_entry,_count from counters where machine_id=_machine_id and mitarbeiter_id=_mitarbeiter_id and auftrags_id=_auftrags_id order by created_at desc limit 1;
   INSERT INTO counters (machine_id,auftrags_id, mitarbeiter_id,count)
VALUES (_machine_id,_auftrags_id,_mitarbeiter_id,_count+1);
   select id,created_at,count into _id,new_entry from counters where machine_id=_machine_id and mitarbeiter_id=_mitarbeiter_id and auftrags_id=_auftrags_id order by created_at desc limit 1;
 select extract(epoch from (new_entry-last_entry)into _diff);
 raise notice '_diff:%',_diff;
 UPDATE counters
SET diff = _diff
WHERE
	id=_id;
 return _count+1;
end;
$$;
