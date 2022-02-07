import HostsPage from  '../pageobjects/hosts.page';
import path = require('path');

describe('SSHMan Hosts', () => {
    it('should be able to add host', async () => {
        await HostsPage.open();
        await HostsPage.addUser({
            key: path.resolve(__dirname, '../../../tests/docker/id_rsa.pub'),
            hostname: 'localhost:2201',
            user: 'root',
            alias: 'host1',
        });
    });

    it('should be able to edit host', async () => {
        await HostsPage.open();
        await HostsPage.editHost('host1', {
            key: path.resolve(__dirname, '../../../tests/docker/id_rsa.pub'),
            hostname: 'localhost:2202',
            user: 'root',
            alias: 'host2',
        });
    });

    it('should be able to delete host', async () => {
        await HostsPage.open();
        await HostsPage.deleteHost('host2');
    });

});

