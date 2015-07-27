namespace go generated
namespace php generated

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

struct Track {
    1: Visit visit,
    2: list<Flash> flashes,
}

struct Mark {
    1: string markID,
    2: string clientID,
    3: map<string, string> fields,
    4: map<string, list<string>> fieldLists,
}

service VisitService {
    Visit findByID(1:string visitID),
    list<Visit> findAll(1:i64 limit, 2:i64 offset),
    list<Visit> findAllBySessionID(1:string sessionID, 2:i64 limit, 3:i64 offset),
    list<Visit> findAllByClientID(1:string clientID, 2:i64 limit, 3:i64 offset),
}

service EventService {
    Event findByID(1:string eventID),
    list<Event> findAll(1:i64 limit, 2:i64 offset),
    Event insert(1: bool enabled, 2: map<string, string> fields, 3: map<string, string> filters),
    Event update(1:Event event),
}

service FlashService {
    Flash findByID(1:string flashID),
    list<Flash> findAll(1:i64 limit, 2:i64 offset),
    list<Flash> findAllByVisitID(1:string visitID),
    list<Flash> findAllByEventID(1:string eventID, 2:i64 limit, 3:i64 offset),
}

service TrackService {
    Track track(1:string sessionID, 2:string clientID, 3:map<string, string> fields),
}

service MarkService {
    Mark findByID(1:string markID),
    Mark findByClientID(1:string clientID),
    list<Mark> findAll(1:i64 limit, 2:i64 offset),
    Mark update(1:Mark mark),
}
