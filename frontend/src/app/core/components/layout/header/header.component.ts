import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';

@Component({
    standalone: true,
    selector: 'app-header',
    templateUrl: 'header.component.html',
    imports: [RouterOutlet],
    styleUrls: ['header.component.scss']
})
export class HeaderComponent {}
