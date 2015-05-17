
data_path = "."

modules_enabled = {
	-- Generally required
	"roster"; -- Allow users to have a roster. Recommended ;)
	"saslauth"; -- Authentication for clients and servers. Recommended if you want to log in.
	"tls"; -- Add support for secure TLS on c2s/s2s connections

	-- Nice to have
	"legacyauth"; -- Legacy authentication. Only used by some old clients and bots.
	"version"; -- Replies to server version requests
	"uptime"; -- Report how long server has been running
	"time"; -- Let others know the time here on this server
	"ping"; -- Replies to XMPP pings with pongs
	"register"; -- Allow users to register on this server using a client and change passwords
};

log = {{to = "console", levels = { min =  "debug" }}}

allow_registration = false;

-- ssl = {
-- 	key = "certs/localhost.key";
-- 	certificate = "certs/localhost.crt";
-- }

allow_unencrypted_plain_auth = true -- god have mercy our client doesn't support anything better yet.  removing this is an objective.

VirtualHost "localhost"

-- Set up a MUC (multi-user chat) room server on conference.example.com:
Component "conference.example.com" "muc"
