let ACL = ./ACL.dhall

let prelude = ./Prelude.dhall

let Map = prelude.Map

in  { Type =
        { groups : Map.Type Text (List Text)
        , hosts : Map.Type Text Text
        , tagOwners : Map.Type Text (List Text)
        , acls : List ACL.Type
        }
    , default =
      { groups = Map.empty Text (List Text)
      , hosts = Map.empty Text Text
      , tagOwners = Map.empty Text (List Text)
      , acls = [] : List ACL.Type
      }
    }
