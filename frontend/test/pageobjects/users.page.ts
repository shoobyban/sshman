import CrudPage from './crud.page';

class UsersPage extends CrudPage {

    public async addUser (user: object): Promise<void> {
        await this.addItem('email', user);
    }

    public async editUser (email: string, user: object): Promise<void> {
        await this.editItem(email, 'email', user);
    }

    public async deleteUser (email: string): Promise<void> {
        await this.deleteItem(email);
    }

    public open(): Promise<string> {
        return super.open('users');
    }
}

export default new UsersPage();
