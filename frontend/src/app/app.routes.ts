import {Routes} from '@angular/router';

import {BingoPageComponent} from './pages/bingo/bingo.page.component';
import {LoginPageComponent} from './pages/login/login.page.component';

export const appRoutes: Routes = [
    {
        path: '',
        component: BingoPageComponent
    },
    {
        path: 'login',
        component: LoginPageComponent
    }
];
