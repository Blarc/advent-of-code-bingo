import {Routes} from '@angular/router';

import {HeaderComponent} from './core/components/layout/header/header.component';
import {AboutPageComponent} from './pages/about/about.page.component';
import {BingoPageComponent} from './pages/bingo/bingo.page.component';
import {LoginPageComponent} from './pages/login/login.page.component';
import {ThreadsPageComponent} from './pages/threads/threads.page.component';

export const appRoutes: Routes = [
    {
        path: '',
        component: HeaderComponent,
        children: [
            {
                path: '',
                component: BingoPageComponent
            },
            {
                path: 'login',
                component: LoginPageComponent
            },
            {
                path: 'threads',
                component: ThreadsPageComponent
            },
            {
                path: 'about',
                component: AboutPageComponent
            }
        ]
    }
];
