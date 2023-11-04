import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';

@Component({
    standalone: true,
    selector: 'app-root',
    templateUrl: './app.component.html',
    imports: [RouterOutlet],
    styleUrls: ['./app.component.scss']
})
export class AppComponent {}
