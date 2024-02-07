
#!/bin/sh

# You could probably do this fancier and have an array of extensions
# to create, but this is mostly an illustration of what can be done

psql -v ON_ERROR_STOP=1 --username "api" <<EOF
create extension ulid;
select * FROM pg_extension;
EOF