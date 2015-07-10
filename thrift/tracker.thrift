namespace go tracker

struct Visit {
  1: string visitID,
  2: string sessionID,
  3: string clientID,
  4: i64 timestamp,
  5: map<string, string> fields,
}

struct Event {
  1: string eventID,
  2: bool enabled,
  3: map<string, string> fields,
  4: map<string, string> filters,
}

struct Flash {
  1: string flashID,
  2: string visitID,
  3: string eventID,
  4: i64 timestamp,
  5: map<string, string> visitFields,
  6: map<string, string> eventFields,
}

service TrackerService {
   list<Flash> track(1:string sessionID, 2:string clientID, 3:map<string, string> fields);
}

service VisitService {
   list<Visit> findAll(1:i64 limit, 2:i64 offset),

   list<Visit> findAllBySessionID(1:string sessionID),

   list<Visit> findAllByClientID(1:string clientID);
}

service EventService {
   list<Event> findAll(1:i64 limit, 2:i64 offset),

   Event findByID(1:string eventID),

   void insert(1:Event event),

   void update(1:Event event);
}

service FlashService {
   list<Flash> findAll(1:i64 limit, 2:i64 offset),

   list<Flash> findAllByVisitID(1:string visitID),

   list<Flash> findAllByEventID(1:string eventID, 2:i64 limit, 3:i64 offset),

   Flash findByID(1:string eventID),
}