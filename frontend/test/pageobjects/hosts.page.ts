import CrudPage from './crud.page';

class HostsPage extends CrudPage {

    public async addUser (user: object): Promise<void> {
        await this.addItem('alias', user);
    }

    public async editHost (alias: string, host: object): Promise<void> {
        await this.editItem(alias, 'alias', host);
    }

    public async deleteHost (alias: string): Promise<void> {
        await this.deleteItem(alias);
    }

    public open(): Promise<string> {
        return super.open('hosts');
    }
}

export default new HostsPage();
