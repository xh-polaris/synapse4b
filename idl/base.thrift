namespace go base

struct App {
    1: string name = "",
}

struct Response {
    1:          i32                code  = 0 ,
    2:          string             msg   = "",
    3: optional map<string,string> extra     (api.body="extra,omitempty"),
}

struct Page {
    1: optional i32    page   = 1,
    2: optional i32    size   = 10,
    3: optional string cursor = "",
}