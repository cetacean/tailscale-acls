let ts = ./types/package.dhall

let ACL = ts.ACL

let Policy = ts.Policy

let Map = ts.Prelude.Map

in  Policy::{
    , groups =
      [ { mapKey = "group:admins"
        , mapValue = [ "xe@github", "victorvalenca@github" ]
        }
      , { mapKey = "group:within"
        , mapValue =
          [ "xe@github"
          , "victorvalenca@github"
          , "twi@github"
          , "heartmender@github"
          ]
        }
      ]
    , hosts = toMap { vm-subnet = "10.77.0.0/17" }
    , tagOwners =
      [ { mapKey = "tag:nixos", mapValue = [ "group:admins" ] }
      , { mapKey = "tag:hypervisor", mapValue = [ "group:admins" ] }
      , { mapKey = "tag:hetzner", mapValue = [ "group:admins" ] }
      , { mapKey = "tag:ci", mapValue = [ "group:admins" ] }
      , { mapKey = "tag:alpine", mapValue = [ "group:admins" ] }
      , { mapKey = "tag:vm", mapValue = [ "group:admins" ] }
      , { mapKey = "tag:alrest", mapValue = [ "group:admins" ] }
      , { mapKey = "tag:sensitive", mapValue = [ "group:admins" ] }
      ]
    , acls =
      [ ACL::{ src = [ "*" ], dst = [ "vm-subnet:*" ] }
      , ACL::{ src = [ "*" ], dst = [ "autogroup:self:*" ] }
      , ACL::{ src = [ "group:admins" ], dst = [ "*:*" ] }
      , ACL::{ src = [ "*" ], dst = [ "*:*" ] }
      ]
    }
