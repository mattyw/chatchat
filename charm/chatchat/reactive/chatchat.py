import gobinary

from charmhelpers.core import hookenv
from charms.reactive import when


@when('gobinary.started')
def simple_server_start():
    hookenv.open_port(2000)
    hookenv.status_set('active', 'Ready')
