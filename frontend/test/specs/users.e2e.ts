import UsersPage from  '../pageobjects/users.page';
import path = require('path');

describe('SSHMan Users', () => {
    it('should be able to add user', async () => {
        await UsersPage.open();
        await UsersPage.addUser({
            keyfile: path.resolve(__dirname, '../../../tests/docker/id_rsa.pub'),
            email: 'sam@host1',
        });
    });

    it('should be able to edit user', async () => {
        await UsersPage.open();
        await UsersPage.editUser('sam@host1', {
            keyfile: path.resolve(__dirname, '../../../tests/docker/id_rsa.pub'),
            email: 'sam@host2',
        });
    });
});
